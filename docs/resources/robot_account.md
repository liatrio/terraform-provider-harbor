# Robot Account Resource

Manages a robot account within a Harbor project.

## Example Usage

```hcl
resource "harbor_project" "example" {
  name = "example"
}

resource "harbor_robot_account" "example" {
  project_id = harbor_project.example.id

  name = "robot$example"
  access {
    resource = "image"
    action = "pull"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The object ID of the Harbor project this robot account
will be created under. Changing this forces a new resource to be created.
* `name` - (Required) The name of the robot account, which must begin with 'robot$'.
Changing this forces a new resource to be created.
* `access` - (Required) A block determining the access a robot account is granted.
Changing this forces a new resource to be created.
  * `resource` - (Required) Denotes the resource that access is granted for.
Supported values are `image` and `helm-chart`.
  * `action` - (Required) Denotes the action the robot account will be able to
preform on the corrosponding resource. Supported vaules are `push` and `pull`
* `description` - (Optional) A description of this robot account.
Changing this forces a new resource to be created.
* `disabled` - (Optional) If `true` this robot account is disabled and can not
be used. Defaults to `false`
* `expires_at` - (Optional) Denotes the date and time at which the robot account's
authentication token will expire. Set with an RFC3339 UTC formatted string. If
`expires_at` isn't set, the authentication token will never expire.

## Attribute Reference

The following attributes are exported:

* `id` - The object ID of the robot account.
* `token` - The authentication token for the robot account.
