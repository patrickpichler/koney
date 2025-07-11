# 🔔 Alert Export

When hackers access traps, Koney logs an event in the `alerts` container. You can also configure Koney to export these alerts to external systems. For this, you simply add an `DeceptionAlertExport` resource to your cluster.

ℹ️ **Note**: All `DeceptionAlertExport` resources and referenced `Secret` resources must be in the `koney-system` namespace. Resources in other namespaces will be ignored.

At the moment, we support exporting alerts to the following systems:

- [Dynatrace Security Events](#dynatrace-security-events)

## Dynatrace Security Events

1. Open the **'Access Tokens'** app in your Dynatrace environment
2. Create a new access token with scope `openpipeline.events_security`
3. Store the API token in a `Secret` resource in the `koney-system` namespace

```yaml
kubectl create secret generic -n koney-system dynatrace-api-token \
  --from-literal=apiToken=dt0c01.EXAMPLE_TOKEN_REPLACE_THE_ENTIRE_STRING \
  --from-literal=apiUrl=https://ENVIRONMENTID.live.dynatrace.com
```

4. Create a `DeceptionAlertExport` resource in your cluster and point it to the `dnyatrace-api-token` secret:

```yaml
apiVersion: research.dynatrace.com/v1alpha1
kind: DeceptionAlertExport
metadata:
  name: deceptionalertexport-dynatrace
  namespace: koney-system
spec:
  dynatrace:
    secretName: dynatrace-api-token
    severity: HIGH
```

The `dynatrace` section contains the following fields:

- `secretName`: The name of the `Secret` resource containing the `apiToken` and `apiUrl` fields.
- `severity`: The severity of the alert upon ingest. Possible values are `CRITICAL`, `HIGH`, `MEDIUM`, and `LOW`. The default value is `HIGH`.

To apply a deception alert export resource, use the following command:

```sh
kubectl apply -f <deceptionalertexport-file>.yaml
```
