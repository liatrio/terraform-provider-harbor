package harbor

import (
	"fmt"
)

type Webhook struct {
	UpdateTime   string           `json:"update_time,omitempty"`
	Description  string           `json:"description,omitempty"`
	Creator      string           `json:"creator,omitempty"`
	CreationTime string           `json:"creation_time,omitempty"`
	Enabled      bool             `json:"enabled"`
	EventTypes   []string         `json:"event_types"`
	ProjectID    int              `json:"project_id"`
	Id           int64            `json:"id"`
	Name         string           `json:"name,omitempty"`
	Metadata     WebhookTargetObj `json:"metadata"`
}

type WebhookTargetObj struct {
	Type           string `json:"type"`
	AuthHeader     string `json:"auth_header"`
	SkipCertVerify bool   `json:"skip_cert_verify"`
	Address        string `json:"address"`
}

func (client *Client) GetWebhook(id string) (*Webhook, error) {
	var webhook *Webhook

	err := client.get(APIURLVersion2, id, &webhook, nil)
	if err != nil {
		return nil, err
	}

	return webhook, nil
}

func (client *Client) NewWebhook(projectID string, webhook *Webhook) (string, error) {
	_, location, err := client.post(APIURLVersion2, fmt.Sprintf("%s/webhook/policies", projectID), webhook)
	if err != nil {
		return "", err
	}

	return location, nil
}

func (client *Client) UpdateWebhook(id string, webhook *Webhook) error {
	return client.put(APIURLVersion2, id, webhook)
}

func (client *Client) DeleteWebhook(id string) error {
	return client.delete(APIURLVersion2, id, nil)
}
