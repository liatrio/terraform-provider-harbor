package harbor

import (
	"crypto/tls"
	"net/http"
)

type Client struct {
	baseUrl    string
	username   string
	password   string
	httpClient *http.Client
}

func NewClient(baseUrl string, username string, password string, tlsInsecureSkipVerify bool) (*Client, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: tlsInsecureSkipVerify},
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	client := &Client{
		baseUrl:    baseUrl,
		username:   username,
		password:   password,
		httpClient: httpClient,
	}

	return client, nil
}
