terraform {
  required_providers {
    todoist = {
      source = "github.com/andreaswwilson/todoist"
    }
  }
}

provider "todoist" {
}

data "todoist_projects" "example" {
  id = "2328048508"
}
