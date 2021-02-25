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

resource "harbor_webhook" "webhook" {
  name = "webhook"
  project_id = harbor_project.project.id
  event_types = ["DELETE_ARTIFACT"]
  target {
    type = "http"
    address = "http://domain.example/webhook"
  }
}

resource "harbor_label" "label" {
  name = "label"
  color = "#FFFFFF"
  description = "An example label"
  project_id = harbor_project.project.id
}
