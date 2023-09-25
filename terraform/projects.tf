resource "gitlab_project" "project" {
  for_each = var.repos

  name        = each.key
  namespace_id =  gitlab_group.services.id  
}


resource "gitlab_repository_file" "hello_world" {
  for_each = gitlab_project.project

  project        = each.value.id
  file_path      = "hello_world.txt"
  branch         = "main"
  content        = base64encode(file("./hello_world.txt"))
  author_email   = "terraform@keyproland.home"
  author_name    = "Terraform"
  commit_message = "say hello world"
}