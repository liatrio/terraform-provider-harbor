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
	Deleted            string  `json:"deleted"`
	OwnerId            int32   `json:"owner_id"`
	RepoCount          int     `json:"repo_count"`
	CreationTime       string  `json:"creation_time"`
	Togglable          bool    `json:"togglable"`
	ProjectId          int32   `json:"project_id"`
	CurrentUserRoleIds []int32 `json:"current_user_role_ids"`
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

func (client *Client) GetProject(id string) (*Project, error) {
	var project *Project

	err := client.get(fmt.Sprintf("projects/%s", id), &project, nil)
	if err != nil {
		return nil, err
	}

	return project, nil
}
