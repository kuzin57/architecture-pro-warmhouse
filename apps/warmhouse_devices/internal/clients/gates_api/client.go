package gatesapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/warmhouse/warmhouse_devices/internal/config"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(conf *config.Config) *Client {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &Client{httpClient: httpClient}
}

func (c *Client) ActivateGate(ctx context.Context, host string) error {
	url := fmt.Sprintf("http://%s/activate", host)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to activate gate: %s", resp.Status)
	}

	return nil
}

func (c *Client) DeactivateGate(ctx context.Context, host string) error {
	url := fmt.Sprintf("http://%s/deactivate", host)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to deactivate gate: %s", resp.Status)
	}

	return nil
}

func (c *Client) GetGateStatus(ctx context.Context, host string) (*StatusResponse, error) {
	url := fmt.Sprintf("http://%s/status", host)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var status StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &status, nil
}

func (c *Client) ChangeGateState(ctx context.Context, host string, state GateState) error {
	url := fmt.Sprintf("http://%s/changestate", host)

	body, err := json.Marshal(ChangeStateRequest{State: state})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to change gate state: %s", resp.Status)
	}

	return nil
}

func (c *Client) HealthCheck(ctx context.Context, host string) error {
	url := fmt.Sprintf("http://%s/health", host)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to health check: %s", resp.Status)
	}

	return nil
}
