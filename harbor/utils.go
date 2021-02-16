package harbor

func (client *Client) GetResource(id string) (interface{}, error) {
	resource := make(map[string]interface{})

	err := client.get(APIURLVersion2, id, &resource, nil)
	if err != nil {
		return nil, err
	}

	return resource, nil
}
