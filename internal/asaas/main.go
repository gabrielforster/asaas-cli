package asaas

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Webhook struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	URL          string   `json:"url"`
	Email        string   `json:"email"`
	Enabled      bool     `json:"enabled"`
	Interrupted  bool     `json:"interrupted"`
	APIVersion   int      `json:"apiVersion"`
	HasAuthToken bool     `json:"hasAuthToken"`
	SendType     string   `json:"sendType"`
	Events       []string `json:"events"`
}

type WebhookListResponse struct {
	Object     string    `json:"object"`
	HasMore    bool      `json:"hasMore"`
	TotalCount int       `json:"totalCount"`
	Limit      int       `json:"limit"`
	Offset     int       `json:"offset"`
	Data       []Webhook `json:"data"`
}

type ClientConfig struct {
	BaseURL    string
	HTTPClient *http.Client
	APIKey     string
}

type AsaasClient interface {
	ListWebhooks() ([]Webhook, error)
	UpdateWebhookURL(id string, webhook Webhook) (*Webhook, error)
}

type AsaasWebhookClient struct {
	config ClientConfig
}

func NewAsaasWebhookClient(config ClientConfig) *AsaasWebhookClient {
	return &AsaasWebhookClient{
		config,
	}
}

func (c *AsaasWebhookClient) ListWebhooks() ([]Webhook, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		c.config.BaseURL+"/webhooks",
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("access_token", c.config.APIKey)

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list webhooks: %s", resp.Status)
	}

	var response WebhookListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Data, nil
}

