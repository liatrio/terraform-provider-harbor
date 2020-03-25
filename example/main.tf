provider "harbor" {
  url = "http://localhost:30002"
  username = ""
  password = ""
}

resource "harbor_project" "example" {
  name = "example"
}
