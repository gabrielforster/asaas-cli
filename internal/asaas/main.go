package asaas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	ToggleWebhookSync(id string, enabled bool) (*Webhook, error)
}

type AsaasWebhookClient struct {
	config ClientConfig
}

func NewAsaasWebhookClient(config ClientConfig) *AsaasWebhookClient {
	return &AsaasWebhookClient{
		config,
	}
}

func (c *AsaasWebhookClient) createRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.config.BaseURL+endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("access_token", c.config.APIKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *AsaasWebhookClient) executeRequest(req *http.Request, expectedStatus int) (*http.Response, error) {
	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	if resp.StatusCode != expectedStatus {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	return resp, nil
}

func (c *AsaasWebhookClient) decodeResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

func (c *AsaasWebhookClient) makeJSONRequest(method, endpoint string, body interface{}, expectedStatus int) (*http.Response, error) {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := c.createRequest(method, endpoint, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return nil, err
	}

	return c.executeRequest(req, expectedStatus)
}

func (c *AsaasWebhookClient) ListWebhooks() ([]Webhook, error) {
	req, err := c.createRequest(http.MethodGet, "/webhooks", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.executeRequest(req, http.StatusOK)
	if err != nil {
		return nil, fmt.Errorf("failed to list webhooks: %w", err)
	}

	var response WebhookListResponse
	if err := c.decodeResponse(resp, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (c *AsaasWebhookClient) UpdateWebhookURL(id, newURL string) (*Webhook, error) {
	body := map[string]string{
		"url": newURL,
	}

	resp, err := c.makeJSONRequest(http.MethodPut, "/webhooks/"+id, body, http.StatusOK)
	if err != nil {
		return nil, fmt.Errorf("failed to update webhook %s with new url %s: %w", id, newURL, err)
	}

	var updatedWebhook Webhook
	if err := c.decodeResponse(resp, &updatedWebhook); err != nil {
		return nil, err
	}

	return &updatedWebhook, nil
}

func (c *AsaasWebhookClient) ToggleWebhookSync(id string, enabled bool) (*Webhook, error) {
	body := map[string]bool{
		"interrupted": !enabled,
	}

	resp, err := c.makeJSONRequest(http.MethodPut, "/webhooks/"+id, body, http.StatusOK)
	if err != nil {
		return nil, fmt.Errorf("failed to update webhook %s to enabled %v: %w", id, enabled, err)
	}

	var updatedWebhook Webhook
	if err := c.decodeResponse(resp, &updatedWebhook); err != nil {
		return nil, err
	}

	return &updatedWebhook, nil
}
