provider "harbor" {
  harbor_url = "harbor.example.com"
}

resource "harbor_project" "project" {
  name = "parker_lab"
}

resource "harbor_robot_account" "robot_account" {
  project_id = harbor_project.project.id
  name = "robot_account"
}
