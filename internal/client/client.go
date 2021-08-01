// Package client provides a Discord client specifically for Discord Interactions API calls outward (not Interactions handling itself).
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type ClientConfig struct {
	BotToken          string
	ClientCredentials string
	ApplicationID     string
	UserAgent         string
	APIRoot           string
}

// GetAuthHeader returns the header to use with Discord API calls.
// If BotToken is set, it will use Bot ${BotToken}, or Bearer ${ClientCredentials} if the opposite is set.
func (cc ClientConfig) GetAuthHeader() (string, error) {
	if cc.BotToken != "" {
		return fmt.Sprintf("Bot %s", cc.BotToken), nil
	}

	if cc.ClientCredentials != "" {
		return fmt.Sprintf("Bearer %s", cc.ClientCredentials), nil
	}

	return "", fmt.Errorf("neither BotToken nor ClientCredentials is set")
}

type InteractionsClient struct {
	config     ClientConfig
	authHeader string
	httpClient *http.Client
}

func NewInteractionsClient(config ClientConfig) (*InteractionsClient, error) {
	authHeader, err := config.GetAuthHeader()
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}

	if config.APIRoot == "" {
		return nil, fmt.Errorf("APIRoot not set")
	}

	if config.ApplicationID == "" {
		return nil, fmt.Errorf("ApplicationID not set")
	}

	if config.UserAgent == "" {
		config.UserAgent = "(+https://github.com/roleypoly/terraform-provider-discord-interactions)"
	}

	return &InteractionsClient{
		config:     config,
		authHeader: authHeader,
		httpClient: httpClient,
	}, nil
}

func (i *InteractionsClient) makeRequest(method string, path string, body interface{}) (*http.Response, error) {
	url, err := url.Parse(fmt.Sprintf("%s/applications/%s%s", i.config.APIRoot, i.config.ApplicationID, path))
	if err != nil {
		return nil, err
	}

	request := &http.Request{
		Method: method,
		URL:    url,
		Header: map[string][]string{
			"authorization": {i.authHeader},
			"user-agent":    {i.config.UserAgent},
			"content-type":  {"application/json"},
		},
	}

	if body != nil {
		bodyBuffer := bytes.Buffer{}

		err = json.NewEncoder(&bodyBuffer).Encode(body)
		if err != nil {
			return nil, err
		}

		request.Body = ioutil.NopCloser(&bodyBuffer)
	}

	response, err := i.httpClient.Do(request)
	return response, err
}

func (i *InteractionsClient) ErrFromResponse(response *http.Response, any ...interface{}) error {
	return fmt.Errorf("InteractionsClient error for request: %s %s, response:\n  code: %d,\n  extra: %v", response.Request.Method, response.Request.URL.String(), response.StatusCode, any)
}
