package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xanzy/go-gitlab"
)

type GitLabClient struct {
	client *gitlab.Client
}

func NewGitLabClient(token string) *GitLabClient {
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL("https://gitlab.keyproland.home/api/v4"))
	if err != nil {
		log.Fatal(err)
	}
	/*
		took me 1h to understand the line bellow
	*/
	return &GitLabClient{client}
}

type ListOptions struct {
	Page    int `url:"page,omitempty" json:"page,omitempty"`
	PerPage int `url:"per_page,omitempty" json:"per_page,omitempty"`
}
type ListGroupProjectsOptions struct {
	ListOptions
	Topic *string `url:"topic,omitempty" json:"topic,omitempty"`
}

func (gc *GitLabClient) ListGroupProjects(groupID string, topic *string, excludeProjectIDs []int) ([]*gitlab.Project, error) {
	// List projects within the group with specified options
	projects, _, err := gc.client.Groups.ListGroupProjects(groupID, &gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
		},
		Topic: topic,
	})
	if err != nil {
		return nil, err
	}

	filteredProjects := []*gitlab.Project{}
	contains := func(excludeProjectIDs []int, id int) bool {
		for _, id_to_exclude := range excludeProjectIDs {
			if id_to_exclude == id {
				return true
			}
		}
		return false
	}
	for _, project := range projects {
		if !contains(excludeProjectIDs, project.ID) {
			filteredProjects = append(filteredProjects, project)
		}
	}
	return filteredProjects, nil
}

func CreateFileIfNotExists() {
	terraformFile := "./terraform/terraform.tfvars"

	file, err := os.Create(terraformFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.Printf("tfvars file created")
}

func main() {
	token := "glpat-GbKUURYm8xVF9TFrPR3B"
	//"os.Getenv("GITLAB_TOKEN")
	if token == "" {
		log.Fatal("GITLAB_TOKEN environment variable not set")
	}

	gitlab := NewGitLabClient(token)

	topic := ""
	excludeProjectIDs := []int{5}
	projects, err := gitlab.ListGroupProjects("12", &topic, excludeProjectIDs)

	CreateFileIfNotExists()

	if err != nil {
		log.Fatalf("Failed to list group projects: %v", err)
	}

	for _, project := range projects {
		fmt.Printf("Project Name: %s\n", project.Name)
	}
}
