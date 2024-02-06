terraform {
  required_providers {
    todoist = {
      source = "github.com/andreaswwilson/todoist"
    }
  }
}

provider "todoist" {
}

resource "todoist_projects" "example" {
  name = "andreas43"
}
