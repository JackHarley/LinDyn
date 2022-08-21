package myip

import (
	"fmt"
	"io"
	"net/http"
)

const IPv4URL = "https://api4.my-ip.io/ip"

type Client interface {
	GetIPv4() (string, error)
}

type client struct {
	httpClient http.Client
}

func NewClient(httpClient http.Client) *client {
	return &client{httpClient}
}

func (c *client) GetIPv4() (string, error) {
	resp, err := c.httpClient.Get(IPv4URL)
	if err != nil {
		return "", fmt.Errorf("failed to get ip from myip: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-ok response from myip: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body from myip: %w", err)
	}

	return string(body), nil
}
