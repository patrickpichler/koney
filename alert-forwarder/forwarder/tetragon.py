# Copyright (c) 2025 Dynatrace LLC
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

import json
import re
from collections import defaultdict

from kubernetes import client

from .fingerprint import (
    KONEY_FINGERPRINT,
    encode_fingerprint_in_cat,
    encode_fingerprint_in_echo,
)
from .types import ContainerMetadata, KoneyAlert, PodMetadata, ProcessMetadata

# group, version, plural of the Tetragon TracingPolicy CRD
TETRAGON_TRACING_POLICIES_GVP = "cilium.io", "v1alpha1", "tracingpolicies"

# the namespace where Tetragon is assumed to be running
TETRAGON_NAMESPACE = "kube-system"
# all tracing policies created by Koney have this prefix
TETRAGON_POLICY_PREFIX = "koney-tracing-policy-"
# the label selector to find Tetragon pods
TETRAGON_POD_LABEL_SELECTOR = "app.kubernetes.io/name=tetragon"
# the container name where Tetragon logs are written
TETRAGON_POD_CONTAINER_NAME = "export-stdout"
# the label key that references the deception policy in a tracing policy
TETRAGON_DECEPTION_POLICY_REF = "koney/deception-policy"

# various error messages
K8S_POLICY_RESOLUTION_ERROR = "failed to resolve DeceptionPolicy name"
K8S_TRAP_TYPE_RESOLUTION_ERROR = "failed to resolve trap type and metadata"

# stores hashes of already processed events to prevent duplicates
event_cache = set()


def read_tetragon_events(since_seconds=60) -> dict[str, list[dict]]:
    v1 = client.CoreV1Api()

    pod_list = v1.list_namespaced_pod(
        namespace=TETRAGON_NAMESPACE,
        label_selector=TETRAGON_POD_LABEL_SELECTOR,
    )

    events_per_policy = defaultdict(list)
    for pod in pod_list.items:
        loglines = v1.read_namespaced_pod_log(
            name=pod.metadata.name,
            namespace=TETRAGON_NAMESPACE,
            container=TETRAGON_POD_CONTAINER_NAME,
            since_seconds=since_seconds,
        )

        for line in loglines.splitlines():
            # quickly filter-out lines that cannot match
            if TETRAGON_POLICY_PREFIX not in line:
                continue

            # events are often duplicated because kprobes can trigger multiple times.
            # as a simple de-duplication strategy, we remove the milliseconds from the timestamp.
            # this filters events that are completely identical and occurred within the same second.
            time_pattern = r'("time":"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})\.\d{9}(Z")'
            line = re.sub(time_pattern, r"\1\2", line)

            # parse and check the referenced policy name
            event = json.loads(line)

            if (policy_name := extract_tracing_policy_name(event)) is not None:
                if not policy_name.startswith(TETRAGON_POLICY_PREFIX):
                    continue

                # avoid duplicates
                event_hash = hash(line)
                if event_hash in event_cache:
                    continue

                event_cache.add(event_hash)
                events_per_policy[policy_name].append(event)

    # returns the list of events (value) grouped by their policy name (key)
    return events_per_policy


def map_tetragon_event(event: dict) -> KoneyAlert:
    deception_policy_name = None
    trap_type = "unknown"
    metadata = dict()

    try:
        # attempt to resolve the DeceptionPolicy name (calls Kubernetes API)
        if tracing_policy_name := extract_tracing_policy_name(event):
            deception_policy_name = resolve_deception_policy_name(tracing_policy_name)
    except client.ApiException:
        pass

    # infer trap type and metadata by inspecting the event
    if "process_kprobe" in event:
        k = event["process_kprobe"]
        file_access_fn = ("security_file_permission", "security_mmap_file")
        if k.get("function_name") in file_access_fn:
            trap_type = "filesystem_honeytoken"
            file_path = k.get("args", [{}])[0].get("file_arg", {}).get("path")
            metadata["file_path"] = file_path

    pod = extract_pod_metadata(event)
    process = extract_process_metadata(event)

    # TODO: emit errors if we fail to resolve fields
    return KoneyAlert(
        timestamp=event["time"],
        deception_policy_name=deception_policy_name,
        trap_type=trap_type,
        metadata=metadata,
        pod=pod,
        process=process,
    )


def is_filtered_event(event: KoneyAlert) -> bool:
    arguments = (event.get("process") or {}).get("arguments")
    if arguments:
        fingerprints = [
            encode_fingerprint_in_echo(KONEY_FINGERPRINT),
            encode_fingerprint_in_cat(KONEY_FINGERPRINT),
        ]

        # if any fingerprint is present, filter this event
        return any(fp in arguments for fp in fingerprints)

    # assume not filtered
    return False


def resolve_deception_policy_name(tracing_policy_name: str) -> str:
    api = client.CustomObjectsApi()
    tp = api.get_cluster_custom_object(
        *TETRAGON_TRACING_POLICIES_GVP, tracing_policy_name
    )

    # TODO: fix typing warnings
    return tp["metadata"]["labels"][TETRAGON_DECEPTION_POLICY_REF]


def extract_tracing_policy_name(event: dict) -> str | None:
    # keys might be process_kprobe, process_uprobe, ...
    for key in event.keys():
        if "policy_name" in event[key]:
            return event[key]["policy_name"]


def extract_pod_metadata(event: dict) -> PodMetadata | None:
    # keys might be process_kprobe, process_uprobe, ...
    for key in event.keys():
        if "process" not in event[key]:
            continue
        if "pod" not in event[key]["process"]:
            continue

        p = event[key]["process"]["pod"]
        return PodMetadata(
            name=p.get("name"),
            namespace=p.get("namespace"),
            container=ContainerMetadata(
                id=p.get("container", {}).get("id"),
                name=p.get("container", {}).get("name"),
            ),
        )


def extract_process_metadata(event: dict) -> ProcessMetadata | None:
    # keys might be process_kprobe, process_uprobe, ...
    for key in event.keys():
        if "process" not in event[key]:
            continue

        p = event[key]["process"]
        return ProcessMetadata(
            pid=p.get("pid"),
            cwd=p.get("cwd"),
            binary=p.get("binary"),
            arguments=p.get("arguments"),
        )
