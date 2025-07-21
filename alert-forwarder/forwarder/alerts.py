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
from hashlib import md5
from pathlib import Path

from .types import DynatraceSeverity, KoneyAlert


def create_alert_id(koney_alert: KoneyAlert) -> str:
    koney_alert_str = json.dumps(koney_alert)
    event_hash = md5(koney_alert_str.encode("utf-8"))
    return event_hash.hexdigest().upper()


def create_alert_description(koney_alert: KoneyAlert) -> str:
    if koney_alert["trap_type"] == "filesystem_honeytoken":
        file_path = koney_alert.get("metadata", {}).get("file_path", "?")
        namespace = (koney_alert.get("pod", {}) or {}).get("namespace")
        pod = (koney_alert.get("pod", {}) or {}).get("name")
        namespaced_pod_name = f"{namespace}/{pod}" if namespace and pod else "?"
        return f"Access to honeytoken ({file_path}) in pod ({namespaced_pod_name}) detected"

    return "Koney alert triggered"


def map_to_dynatrace_event(
    koney_alert: KoneyAlert, severity: DynatraceSeverity
) -> dict:
    # create ids and descriptions
    alert_id = create_alert_id(koney_alert)
    alert_description = create_alert_description(koney_alert)

    # resolve fields, or prepare empty defaults
    pod_dict = koney_alert.get("pod", {}) or {}
    node_dict = koney_alert.get("node", {}) or {}
    process_dict = koney_alert.get("process", {}) or {}

    # split process binary into name and path (with pathlib)
    process_binary = Path(process_dict.get("binary", ""))
    process_binary_name = process_binary.name
    process_binary_path = str(process_binary.parent)

    # map severity text to risk score
    risk_score = 0
    severity_norm = severity.lower()
    if severity_norm == "low":
        risk_score = 3.9
    elif severity_norm == "medium":
        risk_score = 6.9
    elif severity_norm == "high":
        risk_score = 8.9
    elif severity_norm == "critical":
        risk_score = 10.0

    # TODO: bump event.version
    payload = {
        "timestamp": koney_alert["timestamp"],
        # koney metadata (flattened)
        "koney.deception_policy_name": koney_alert["deception_policy_name"],
        "koney.trap_type": koney_alert["trap_type"],
        "koney.metadata.file_path": koney_alert.get("metadata", {}).get("file_path"),
        # event metadata
        "event.kind": "SECURITY_EVENT",
        "event.type": "DETECTION_FINDING",
        "event.name": "Detection finding event",
        "event.provider": "Koney",
        "event.version": "2025-07-18",
        "event.id": alert_id,
        "event.description": alert_description,
        # detection metadata
        "detection.type": "KONEY_ALERT",
        # security finding metadata
        "finding.id": alert_id,
        "finding.title": alert_description,
        "finding.description": alert_description,
        "finding.time.created": koney_alert["timestamp"],
        "finding.severity": severity,
        # security event metadata
        "dt.security.risk.level": severity,
        "dt.security.risk.score": risk_score,
        # product metadata
        "product.name": "Koney",
        "product.vendor": "Dynatrace Research",
        # kubernetes metadata
        "k8s.namespace.name": pod_dict.get("namespace"),
        "k8s.node.name": node_dict.get("name"),
        "k8s.pod.name": pod_dict.get("name"),
        "k8s.container.name": pod_dict["container"].get("name"),
        "k8s.container.id": pod_dict["container"].get("id"),
        # process metadata
        "process.executable.name": process_binary_name,
        "process.executable.path": process_binary_path,
        "process.executable.arguments": process_dict.get("arguments"),
        "process.pid": process_dict.get("pid"),
        "process.uid": process_dict.get("uid"),
        "process.cwd": process_dict.get("cwd"),
        # source object metadata (for enrichment)
        "object.type": "KUBERNETES_CONTAINER",
        "object.id": pod_dict["container"].get("id"),
        # not collected by Koney
        # "k8s.cluster.name": "",
        # "k8s.pod.uid": "",
    }

    return payload
