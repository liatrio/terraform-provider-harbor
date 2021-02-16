package harbor

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

func (client *Client) GetProject(id string) (*Project, error) {
	var project *Project

	err := client.get(APIURLVersion2, id, &project, nil)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (client *Client) NewProject(project *ProjectReq) (string, error) {
	_, location, err := client.post(APIURLVersion2, "/projects", project)
	if err != nil {
		return "", err
	}

	return location, nil
}

func (client *Client) UpdateProject(id string, project *ProjectReq) error {
	return client.put(APIURLVersion2, id, project)
}

func (client *Client) DeleteProject(id string) error {
	return client.delete(APIURLVersion2, id, nil)
}
