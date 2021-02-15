package harbor

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Client struct {
	baseURL    string
	username   string
	password   string
	httpClient *http.Client
}

const (
	ApiURLVersion1 = "/api"
	ApiURLVersion2 = "/api/v2.0"
)

func NewClient(baseURL string, username string, password string, tlsInsecureSkipVerify bool) (*Client, error) {
	transport := &http.Transport{
		//nolint:gosec
		TLSClientConfig: &tls.Config{InsecureSkipVerify: tlsInsecureSkipVerify},
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	client := &Client{
		baseURL:    baseURL,
		username:   username,
		password:   password,
		httpClient: httpClient,
	}

	return client, nil
}

func (client *Client) sendRequest(request *http.Request) ([]byte, string, error) {
	request.SetBasicAuth(client.username, client.password)
	request.Header.Add("Content-Type", "application/json")

	requestMethod := request.Method
	requestPath := request.URL.Path

	log.Printf("[DEBUG] Sending %s to %s", requestMethod, requestPath)
	showBody := false
	if request.Body != nil {
		showBody = true
		requestBody, err := request.GetBody()
		if err != nil {
			return nil, "", err
		}

		requestBodyBuffer := new(bytes.Buffer)
		_, err = requestBodyBuffer.ReadFrom(requestBody)
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] Request body: %s", requestBodyBuffer.String())
	}

	dump, err := httputil.DumpRequest(request, showBody)
	if err != nil {
		return nil, "", err
	}
	log.Printf("[DEBUG] %s", dump)

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}

	if response.StatusCode >= 400 {
		return nil, "", &APIError{
			Code:    response.StatusCode,
			Message: fmt.Sprintf("error sending %s request to %s: %s", request.Method, request.URL.Path, response.Status),
		}
	}

	return body, response.Header.Get("Location"), nil
}

func (client *Client) get(apiURL string, path string, resource interface{}, params map[string]string) error {
	resourceURL := client.baseURL + apiURL + path

	request, err := http.NewRequest(http.MethodGet, resourceURL, nil)
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

	body, _, err := client.sendRequest(request)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, resource)
}

func (client *Client) post(apiURL string, path string, requestBody interface{}) ([]byte, string, error) {
	resourceURL := client.baseURL + apiURL + path

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return nil, "", err
	}

	request, err := http.NewRequest(http.MethodPost, resourceURL, bytes.NewReader(payload))
	if err != nil {
		return nil, "", err
	}

	body, location, err := client.sendRequest(request)
	location = strings.Replace(location, apiURL, "", 1)

	return body, location, err
}

func (client *Client) put(apiURL string, path string, requestBody interface{}) error {
	resourceURL := client.baseURL + apiURL + path

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPut, resourceURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	_, _, err = client.sendRequest(request)

	return err
}

func (client *Client) delete(apiURL string, path string, requestBody interface{}) error {
	resourceURL := client.baseURL + apiURL + path

	var body io.Reader

	if requestBody != nil {
		payload, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}
		body = bytes.NewReader(payload)
	}

	request, err := http.NewRequest(http.MethodDelete, resourceURL, body)
	if err != nil {
		return err
	}

	_, _, err = client.sendRequest(request)

	return err
}
