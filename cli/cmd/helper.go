package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/leejss/simple-json-server/cli/jira"
)

func buildSearchRequest(jql string, startAt, maxResults int, fields []string) jira.SearchRequest {
	return jira.SearchRequest{
		JQL:        jql,
		StartAt:    startAt,
		MaxResults: maxResults,
		Fields:     fields,
	}
}

func buildHTTPRequest(ctx context.Context, baseURL, apiToken string, reqBody jira.SearchRequest) (*http.Request, error) {
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create http request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/rest/api/2/search", bytes.NewBuffer(bodyBytes))

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// 원본 바이트와 파싱된 응답을 함께 리턴. client, req => doRequest -> bytes, response, err
func doRequest(client *http.Client, req *http.Request) ([]byte, *jira.SearchResponse, error) {
	resp, err := client.Do(req)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	var parsed jira.SearchResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return body, &parsed, nil

}

func prettyJson(body []byte) ([]byte, error) {
	var prettyJson bytes.Buffer
	if err := json.Indent(&prettyJson, body, "", "  "); err != nil {
		return nil, fmt.Errorf("failed to indent JSON: %w", err)
	}
	return prettyJson.Bytes(), nil
}

func processYear() {}

func paginate() {}
