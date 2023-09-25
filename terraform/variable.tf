variable "gitlab_token" {
  type    = string
  default = "glpat-GbKUURYm8xVF9TFrPR3B"
}

variable "repos" {
  description = "Map of repositories"
  type        = map(object({
    id = number
  }))
}