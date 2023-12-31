package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

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

type File struct {
	Path string
}

type Project struct {
	ID   int
	Name string
}

func (f *File) CreateFileIfNotExists() *File {

	_, err := os.Stat(f.Path)

	if err == nil {
		fmt.Printf("File %s exists!\n", f.Path)
	} else if os.IsNotExist(err) {
		fmt.Printf("File %s does not exist.\n", f.Path)
		file, err := os.Create(f.Path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		log.Printf("tfvars file created")

	} else {
		fmt.Printf("Error checking if file exists: %s", err)
	}
	return f
}

func (f *File) isEmpty() (bool, error) {
	file, err := os.Open(f.Path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return false, err
	}

	return stat.Size() == 0, nil
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

	terraformFile := "./terraform/terraform.tfvars"

	isEmpty, err := (&File{Path: terraformFile}).CreateFileIfNotExists().isEmpty()
	if err != nil {
		log.Fatal(err)
	}

	if isEmpty {
		log.Printf("The file is empty.")
		log.Printf("Updating the tfvars file")

		file, err := os.OpenFile(terraformFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		fmt.Fprintln(file, "repos = {")
		if err != nil {
			log.Fatalf("Error getting current working directory: %s", err)
		}
		for _, project := range projects {
			fmt.Fprintf(file, `  "%s" = {
    id = %d
  }
`, project.Name, project.ID)
		}
		fmt.Fprintln(file, "}")

		originalDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting current working directory: %s", err)
		}

		if err := os.Chdir("terraform"); err != nil {
			log.Fatalf("Error changing directory: %s", err)
		}

		for _, project := range projects {
			cmd := exec.Command("terraform", "import",
				fmt.Sprintf("gitlab_project.project[\"%s\"]", project.Name),
				fmt.Sprintf("%d", project.ID))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				log.Fatalf("Error running terraform import: %s", err)
			}
		}

		if err := os.Chdir(originalDir); err != nil {
			log.Fatalf("Error changing back to original directory: %s", err)
		}

		fmt.Println("Updating terraform state completed")
	} else {
		log.Printf("The file is not empty.")
		log.Printf("Cleaning the file before updating it")

	}

	if err != nil {
		log.Fatalf("Failed to list group projects: %v", err)
	}

	for _, project := range projects {
		fmt.Printf("Project Name: %s\n", project.Name)
	}
}
