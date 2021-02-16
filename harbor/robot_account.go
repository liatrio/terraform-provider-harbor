package harbor

import (
	"encoding/json"
	"fmt"
)

type RobotAccountCreate struct {
	Access      []RobotAccountAccess `json:"access,omitempty"`
	Name        string               `json:"name,omitempty"`
	Description string               `json:"description,omitempty"`
	ExpiresAt   int64                `json:"expires_at,omitempty"`
}

type RobotAccountPostRep struct {
	Token string `json:"token"`
	Name  string `json:"name"`
}

type RobotAccountUpdate struct {
	Disabled bool `json:"disabled,omitempty"`
}

type RobotAccount struct {
	Description  string `json:"description"`
	UpdateTime   string `json:"update_time"`
	CreationTime string `json:"creation_time"`
	ExpiresAt    int    `json:"expires_at"`
	Disabled     bool   `json:"disabled"`
	ProjectID    int    `json:"project_id"`
	ID           int    `json:"id"`
	Name         string `json:"name"`
}

type RobotAccountAccess struct {
	Action   string `json:"action,omitempty"`
	Resource string `json:"resource,omitempty"`
}

func (client *Client) GetRobotAccount(id string) (*RobotAccount, error) {
	var robot *RobotAccount

	err := client.get(APIURLVersion2, id, &robot, nil)
	if err != nil {
		return nil, err
	}

	return robot, nil
}

func (client *Client) NewRobotAccount(projectID string, robot *RobotAccountCreate) (*RobotAccountPostRep, string, error) {
	body, location, err := client.post(APIURLVersion2, fmt.Sprintf("%s/robots", projectID), robot)
	if err != nil {
		return nil, "", err
	}

	var response *RobotAccountPostRep

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, "", err
	}

	return response, location, nil
}

func (client *Client) UpdateRobotAccount(id string, robot *RobotAccountUpdate) error {
	return client.put(APIURLVersion2, id, robot)
}

func (client *Client) DeleteRobotAccount(id string) error {
	return client.delete(APIURLVersion2, id, nil)
}
