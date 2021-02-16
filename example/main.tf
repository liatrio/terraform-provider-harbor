terraform {
  required_providers {
    harbor = {
      version = ">=0.0.0"
      source = "terraform.local/liatrio/harbor"
    }
  }
}

provider "harbor" {
// Values should be set here, or in ENV
//  url = "http://localhost"
//  username = ""
//  password = ""
}


resource "harbor_project" "project" {
  name = "project"
}


resource "harbor_robot_account" "robot" {
  name = "robot$robot"
  project_id = harbor_project.project.id
  access {
    resource = "image"
    action = "pull"
  }
}
