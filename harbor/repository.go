package harbor

import (
	"fmt"
)

type Repository struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	ProjectID   int           `json:"project_id"`
	Description string        `json:"description"`
	PullCount   int           `json:"pull_count"`
	StarCount   int           `json:"star_count"`
	TagsCount   int           `json:"tags_count"`
	Labels      []interface{} `json:"labels"`
	// CreationTime time.Time     `json:"creation_time"`
	// UpdateTime   time.Time     `json:"update_time"`
}

func (client *Client) GetRepositories(id string) ([]*Repository, error) {
	var repositories []*Repository
	err := client.get(fmt.Sprintf("/repositories?project_id=%s", id), &repositories, nil)
	if err != nil {
		return nil, err
	}

	return repositories, nil
}

func (client *Client) DeleteRepository(id string) error {
	return client.delete(fmt.Sprintf("/repositories/%s", id), nil)
}

func (client *Client) DeleteRepositories(repos []*Repository) error {
	for _, repo := range repos {
		err := client.DeleteRepository(repo.Name)
		if err != nil {
			return err
		}
	}
	return nil
}
