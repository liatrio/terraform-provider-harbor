# Webhook Resource

Manages a webhook policy within a Harbor project.

## Example Usage

```hcl
resource "harbor_project" "example" {
  name = "example"
}

resource "harbor_webhook" "example" {
  project_id = harbor_project.example.id

  name = "example"
  event_types = ["PULL_ARTIFACT", "PUSH_ARTIFACT"]
  target {
    type = "http"
    address = "http://domain.example/webhook/event"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The object ID of the Harbor project this webhook policy
will be created under. Changing this forces a new resource to be created.
* `name` - (Required) The name of the webhook policy.
* `event_types` - (Required) List of events which will cause the webhook to trigger.
Accepted values are:
  * `DELETE_ARTIFACT`
  * `PULL_ARTIFACT`
  * `PUSH_ARTIFACT`
  * `DELETE_CHART`
  * `DOWNLOAD_CHART`
  * `UPLOAD_CHART`
  * `QUOTA_EXCEED`
  * `QUOTA_WARNING`
  * `REPLICATION`
  * `SCANNING_FAILED`
  * `SCANNING_COMPLETED`
  * `TAG_RETENTION`
* `target` - (Required) Nested block detailing webhook target information.
  * `type` - (Required) The type of webhook payload to send.
  Accepted values are `http` and `slack`.
  * `address` - (Required) The endpoint URL for the webhook payload.
  * `auth_header` - (Optional) The auth header to include with webhook payload.
  * `skip_cert_verify` - (Optional) If `true`, skips verification of endpoint's
  tls certificate. Defaults to `true`.
* `enabled` - (Optional) If `false`, webhooks will not trigger even if
configured events occur. Defaults to `true`
* `description` - (Optional) The description of the webhook policy

## Attribute Reference

The following attributes are exported:

* `id` - The object ID of the webhook policy.
