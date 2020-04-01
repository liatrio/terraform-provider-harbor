package harbor

import (
	"fmt"
)

type ProjectReq struct {
	CountLimit  int64  `json:"count_limit,omitempty"`
	ProjectName string `json:"project_name,omitempty"`
	// CVEWhitelist `json:"cve_whitelist"`
	StorageLimit int64           `json:"storage_limit,omitempty"`
	Metadata     ProjectMetadata `json:"metadata,omitempty"`
}

type Project struct {
	UpdateTime         string  `json:"update_time"`
	OwnerName          string  `json:"owner_name"`
	Name               string  `json:"name"`
	Deleted            bool    `json:"deleted"`
	OwnerID            int32   `json:"owner_id"`
	RepoCount          int     `json:"repo_count"`
	CreationTime       string  `json:"creation_time"`
	Togglable          bool    `json:"togglable"`
	ProjectID          int32   `json:"project_id"`
	CurrentUserRoleIDs []int32 `json:"current_user_role_ids"`
	ChartCount         int     `json:"chart_count"`
	// CVEWhitelist `json:"cve_whitelist"`
	Metadata ProjectMetadata `json:"metadata"`
}

type ProjectMetadata struct {
	EnableContentTrust   string `json:"enable_content_trust,omitempty"`
	AutoScan             string `json:"auto_scan,omitempty"`
	Severity             string `json:"severity,omitempty"`
	ReuseSysCveWhitelist string `json:"reuse_sys_cve_whitelist,omitempty"`
	Public               string `json:"public,omitempty"`
	PreventVul           string `json:"prevent_vul,omitempty"`
}

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

func (client *Client) GetProject(id string) (*Project, error) {
	var project *Project

	err := client.get(id, &project, nil)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (client *Client) NewProject(project *ProjectReq) (string, error) {
	_, location, err := client.post("/projects", project)
	if err != nil {
		return "", err
	}

	return location, nil
}

func (client *Client) UpdateProject(id string, project *ProjectReq) error {
	return client.put(id, project)
}

func (client *Client) GetRepositories(id string) ([]*Repository, error) {
	var repositories []*Repository
	err := client.get(fmt.Sprintf("/repositories?project_id=%s", id), &repositories, nil)
	if err != nil {
		return nil, err
	}

	return repositories, nil
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

func (client *Client) DeleteRepository(id string) error {
	return client.delete(fmt.Sprintf("/repositories/%s", id), nil)
}

func (client *Client) DeleteProject(id string) error {
	return client.delete(id, nil)
}
