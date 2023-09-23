# you still need to import the groups manually, comment projects.tf file
# You can import a group state using terraform import <resource> <id>

resource "gitlab_group" "keyproland" {
  path = "keyproland"
  name = "keyproland"
}

resource "gitlab_group" "france" {
  name      = "france"
  path      = "france"
  parent_id = gitlab_group.keyproland.id
}

resource "gitlab_group" "services" {
  name      = "services"
  path      = "services"
  parent_id = gitlab_group.france.id
}
