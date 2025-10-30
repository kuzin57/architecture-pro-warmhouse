package temperatureapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

func (c *Client) GetTemperature(ctx context.Context, host string, location string) (*TemperatureResponse, error) {
	baseURL := fmt.Sprintf("http://%s/temperature", host)

	params := url.Values{}
	params.Add("location", location)

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get temperature: %s", resp.Status)
	}

	var temperature TemperatureResponse
	if err := json.NewDecoder(resp.Body).Decode(&temperature); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &temperature, nil
}

func (c *Client) GetTemperatureBySensorID(ctx context.Context, host string, sensorID string) (*TemperatureResponse, error) {
	url := fmt.Sprintf("http://%s/temperature/%s", host, sensorID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get temperature by sensor ID: %s", resp.Status)
	}

	var temperature TemperatureResponse
	if err := json.NewDecoder(resp.Body).Decode(&temperature); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &temperature, nil
}

func (c *Client) HealthCheck(ctx context.Context, host string) error {
	url := fmt.Sprintf("http://%s/health", host)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to health check: %s", resp.Status)
	}

	return nil
}
