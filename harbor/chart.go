package harbor

import (
	"fmt"
	"time"
)

type Chart struct {
	Name          string    `json:"name"`
	TotalVersions int       `json:"total_versions"`
	LatestVersion string    `json:"latest_version"`
	Created       time.Time `json:"created"`
	Updated       time.Time `json:"updated"`
	Icon          string    `json:"icon"`
	Home          string    `json:"home"`
	Deprecated    bool      `json:"deprecated"`
}

func (client *Client) GetCharts(id string) ([]*Chart, error) {
	var charts []*Chart
	err := client.get(ApiURLVersion1, fmt.Sprintf("/chartrepo/%s/charts", id), &charts, nil)
	if err != nil {
		return nil, err
	}

	return charts, nil
}

func (client *Client) DeleteChart(project string, chart string) error {
	return client.delete(ApiURLVersion1, fmt.Sprintf("/chartrepo/%s/charts/%s", project, chart), nil)
}

func (client *Client) DeleteCharts(project string, charts []*Chart) error {
	for _, chart := range charts {
		err := client.DeleteChart(project, chart.Name)
		if err != nil {
			return err
		}
	}
	return nil
}
