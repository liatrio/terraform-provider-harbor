package harbor

import (
	"fmt"
	"time"
)

type Repository struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	ProjectID    int           `json:"project_id"`
	Description  string        `json:"description"`
	PullCount    int           `json:"pull_count"`
	StarCount    int           `json:"star_count"`
	TagsCount    int           `json:"tags_count"`
	Labels       []interface{} `json:"labels"`
	CreationTime time.Time     `json:"creation_time"`
	UpdateTime   time.Time     `json:"update_time"`
}

func (client *Client) GetRepositories(projectName string) ([]*Repository, error) {
	var repositories []*Repository
	err := client.get(APIURLVersion2, fmt.Sprintf("/projects/%s/repositories", projectName), &repositories, nil)
	if err != nil {
		return nil, err
	}

	return repositories, nil
}

func (client *Client) DeleteRepository(projectName string, repoName string) error {
	repo := repoName[len(projectName)+1:] // repoName has the format projectName/repo, strip the projectName/ prefix

	return client.delete(APIURLVersion2, fmt.Sprintf("/projects/%s/repositories/%s", projectName, repo), nil)
}

func (client *Client) DeleteRepositories(projectName string, repos []*Repository) error {
	for _, repo := range repos {
		err := client.DeleteRepository(projectName, repo.Name)
		if err != nil {
			return err
		}
	}
	return nil
}
