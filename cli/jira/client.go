package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Issue struct {
	Key    string `json:"key"`
	Fields struct {
		Summary     string    `json:"summary"`
		Created     time.Time `json:"created"`
		Description string    `json:"description"`
	} `json:"fields"`
}

type SearchRequest struct {
	JQL        string   `json:"jql"`
	StartAt    int      `json:"startAt"`
	MaxResults int      `json:"maxResults"`
	Fields     []string `json:"fields"`
}

type SearchResponse struct {
	StartAt    int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
	Total      int     `json:"total"`
	Issues     []Issue `json:"issues"`
}

type JiraClient struct {
	client   *http.Client
	baseURL  string
	apiToken string
}

func NewJiraClient(baseURL, apiToken string) *JiraClient {
	return &JiraClient{
		client:   &http.Client{},
		baseURL:  baseURL,
		apiToken: apiToken,
	}
}

func (c *JiraClient) Search(jql string) ([]byte, error) {

	reqBody := SearchRequest{
		JQL:        jql,
		StartAt:    0,
		MaxResults: 100,
		Fields:     []string{"key", "summary", "created", "description"},
	}

	jsonBody, err := json.Marshal(reqBody)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	reqUrl := fmt.Sprintf("%s/rest/api/2/search", c.baseURL)

	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(jsonBody))

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	return body, nil
}
