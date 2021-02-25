package harbor

type Label struct {
	CreationTime string `json:"creation_time,omitempty"`
	UpdateTime   string `json:"update_time,omitempty"`
	Deleted      bool   `json:"deleted,omitempty"`
	ID           int64  `json:"id,omitempty"`

	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Scope       string `json:"scope"`
	ProjectID   int64  `json:"project_id"`
}

func (client *Client) GetLabel(id string) (*Label, error) {
	var label *Label

	err := client.get(APIURLVersion2, id, &label, nil)
	if err != nil {
		return nil, err
	}

	return label, nil
}

func (client *Client) NewLabel(label *Label) (string, error) {
	_, location, err := client.post(APIURLVersion2, "/labels", label)
	return location, err
}

func (client *Client) UpdateLabel(id string, label *Label) error {
	return client.put(APIURLVersion2, id, label)
}

func (client *Client) DeleteLabel(id string) error {
	return client.delete(APIURLVersion2, id, nil)
}
