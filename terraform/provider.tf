terraform {
  required_providers {
    gitlab = {
      source = "gitlabhq/gitlab"
    }
  }
}

provider "gitlab" {
  token    = "glpat-GbKUURYm8xVF9TFrPR3B"
  base_url = "https://gitlab.keyproland.home/api/v4"
}