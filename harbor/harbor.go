package harbor

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	baseUrl    string
	username   string
	password   string
	httpClient *http.Client
}

const (
	apiUrl = "/api"
)

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

func (client *Client) sendRequest(request *http.Request) ([]byte, error) {

	request.SetBasicAuth(client.username, client.password)
	request.Header.Add("Content-Type", "application/json")

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 400 {
		return nil, errors.New("Bad Request")
	}

	return body, nil
}

func (client *Client) get(path string, resource interface{}, params map[string]string) error {
	resourceUrl := client.baseUrl + apiUrl + path

	request, err := http.NewRequest(http.MethodGet, resourceUrl, nil)
	if err != nil {
		return err
	}

	if params != nil {
		query := url.Values{}
		for k, v := range params {
			query.Add(k, v)
		}
		request.URL.RawQuery = query.Encode()
	}

	body, err := client.sendRequest(request)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, resource)
}

func (client *Client) post(path string, requestBody interface{}) ([]byte, error) {
	resourceUrl := client.baseUrl + apiUrl + path

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, resourceUrl, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	body, err := client.sendRequest(request)

	return body, err
}

func (client *Client) put(path string, requestBody interface{}) error {
	resourceUrl := client.baseUrl + apiUrl + path

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPut, resourceUrl, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	_, err = client.sendRequest(request)

	return err
}

func (client *Client) delete(path string, requestBody interface{}) error {
	resourceUrl := client.baseUrl + apiUrl + path

	var body io.Reader

	if requestBody != nil {
		payload, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}
		body = bytes.NewReader(payload)
	}

	request, err := http.NewRequest(http.MethodDelete, resourceUrl, body)
	if err != nil {
		return err
	}

	_, err = client.sendRequest(request)

	return err
}
