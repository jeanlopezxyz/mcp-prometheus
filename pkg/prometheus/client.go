package prometheus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is an HTTP client for the Prometheus API.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new Prometheus API client.
func NewClient(baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &Client{BaseURL: baseURL, HTTPClient: httpClient}
}

func (c *Client) doGet(path string) ([]byte, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + path)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

// Query executes a PromQL instant query.
func (c *Client) Query(promql string) (string, error) {
	path := "/api/v1/query?query=" + url.QueryEscape(promql)
	body, err := c.doGet(path)
	if err != nil {
		return "", err
	}
	return formatJSON(body)
}

// QueryRange executes a PromQL range query.
func (c *Client) QueryRange(promql string, start, end time.Time, step string) (string, error) {
	params := url.Values{}
	params.Set("query", promql)
	params.Set("start", start.Format(time.RFC3339))
	params.Set("end", end.Format(time.RFC3339))
	params.Set("step", step)
	path := "/api/v1/query_range?" + params.Encode()
	body, err := c.doGet(path)
	if err != nil {
		return "", err
	}
	return formatJSON(body)
}

// GetTargets returns scrape target information.
func (c *Client) GetTargets(state string) (string, error) {
	path := "/api/v1/targets"
	if state != "" {
		path += "?state=" + url.QueryEscape(state)
	}
	body, err := c.doGet(path)
	if err != nil {
		return "", err
	}
	return formatJSON(body)
}

// GetRules returns alerting and recording rules.
func (c *Client) GetRules(ruleType string) (string, error) {
	path := "/api/v1/rules"
	if ruleType != "" {
		path += "?type=" + url.QueryEscape(ruleType)
	}
	body, err := c.doGet(path)
	if err != nil {
		return "", err
	}
	return formatJSON(body)
}

// GetBuildInfo returns Prometheus build information.
func (c *Client) GetBuildInfo() (string, error) {
	body, err := c.doGet("/api/v1/status/buildinfo")
	if err != nil {
		return "", err
	}
	return formatJSON(body)
}

// GetRuntimeInfo returns Prometheus runtime information.
func (c *Client) GetRuntimeInfo() (string, error) {
	body, err := c.doGet("/api/v1/status/runtimeinfo")
	if err != nil {
		return "", err
	}
	return formatJSON(body)
}

func formatJSON(data []byte) (string, error) {
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return string(data), nil
	}
	formatted, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return string(data), nil
	}
	return string(formatted), nil
}
