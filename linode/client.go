package linode

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	APIBaseURL = "https://api.linode.com"

	ErrDomainNotFound = errors.New("failed to find matching domain")
)

type Client interface {
	GetDomainID(domain string) (int, error)
	GetARecords(domainID int, names map[string]struct{}) ([]RetrievedARecord, error)
	UpdateARecord(domainID int, recordID int, update ARecordUpdate) error
}

type client struct {
	httpClient          http.Client
	personalAccessToken string
}

func NewClient(httpClient http.Client, personalAccessToken string) *client {
	return &client{httpClient, personalAccessToken}
}

func (c *client) doGetRequest(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", APIBaseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new http request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+c.personalAccessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send http request: %w", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read http response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-ok http response code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *client) doPutRequest(path string, body []byte) ([]byte, error) {
	reader := bytes.NewReader(body)

	req, err := http.NewRequest("PUT", APIBaseURL+path, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create new http request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+c.personalAccessToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send http request: %w", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read http response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-ok http response code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
