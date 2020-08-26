# Project Resourcee

Manages a project within Harbor

## Example Usage

```hcl
resource "harbor_project" "example" {
  name = "example"
  public = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Harbor project. Changing this forces a new
project to be created.
* `public` - (Optional) If `true` any user will have read permissions to repositories
under this project. Defaults to `false`

## Attribute Reference

The following attributes are exported:

* `id` - The object ID of the Harbor project.
