# Resource: harbor_label

Manages labels within Harbor.

Labels give the ability to mark artifacts that are stored in Harbor.

## Label Scope

Labels can be created at a global scope, which allows them to be
used in any project, or at a project scope, which only allows them
to be used inside the project they're created under.

When creating a label resource, specifying a project ID will create
the label at a project scope. Leaving project ID unspecified will
create the label at a global scope.

## Example Usage

Global Scope

```hcl
resource "harbor_label" "example" {
  name = "example"
  color = "#000000"
  description = "An example label"
}
```

Project Scope

```hcl
resource "harbor_project" "example" {
  name = "example"
}

resource "harbor_label" "example" {
  name = "example"
  color = "#000000"
  description = "An example label"
  project_id = harbor_project.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the label.
* `color` - (Optional) The color of the label, written as a hex color code. (i.e. "#112233")
* `description` - (Optional) An optional description of the label.
* `project_id` - (Optional) The object ID of the Harbor project this label
should be created under. If this value is not set, the label will be created
at a global scope.

## Attribute Reference

The following attributes are exported:

* `id` - The object ID of the Harbor label.
