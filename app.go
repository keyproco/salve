package main

import (
	"fmt"
	"log"

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

func (gc *GitLabClient) ListGroupProjects(groupID string) ([]*gitlab.Project, error) {
	// List projects within the group
	projects, _, err := gc.client.Groups.ListGroupProjects(groupID, nil)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func main() {
	token := "glpat-GbKUURYm8xVF9TFrPR3B"
	//"os.Getenv("GITLAB_TOKEN")
	if token == "" {
		log.Fatal("GITLAB_TOKEN environment variable not set")
	}

	gitlab := NewGitLabClient(token)
	fmt.Println(gitlab)

	projects, err := gitlab.ListGroupProjects("12")
	if err != nil {
		log.Fatalf("Failed to list group projects: %v", err)
	}

	for _, project := range projects {
		fmt.Printf("Project Name: %s\n", project.Name)
	}
}
