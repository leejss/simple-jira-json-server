package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/leejss/simple-json-server/cli/config"
	"github.com/leejss/simple-json-server/cli/internal/storage"
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

	base, err := url.Parse(strings.TrimRight(baseURL, "/"))

	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	searchURL := base.ResolveReference(&url.URL{Path: "/rest/api/2/search"})

	// Create http request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, searchURL.String(), bytes.NewBuffer(bodyBytes))

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

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
		return body, nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var parsed jira.SearchResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return body, nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return body, &parsed, nil

}

func prettyJson(body []byte) ([]byte, error) {
	if len(bytes.TrimSpace(body)) == 0 {
		return nil, fmt.Errorf("empty response body")
	}
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, body, "", "  "); err != nil {
		return nil, fmt.Errorf("failed to indent JSON: %w", err)
	}

	return pretty.Bytes(), nil
}

func processYear(
	ctx context.Context,
	client *http.Client,
	cfg config.Config,
	builder *jira.JQLQueryBuilder,
	year int,
) error {

	fmt.Printf("Processing year: %d\n", year)

	const pageSize = 100

	reqBuilder := func(startAt int) jira.SearchRequest {
		return buildSearchRequest(
			builder.SearchByYear(year, cfg.Username),
			startAt,
			pageSize,
			[]string{"key", "summary", "created", "description"},
		)
	}

	issues, err := paginate(ctx, client, cfg, reqBuilder)
	if err != nil {
		return fmt.Errorf("(%d) paginate: %w", year, err)
	}

	raw, err := json.Marshal(issues)
	if err != nil {
		return fmt.Errorf("(%d) marshal issues: %w", year, err)
	}

	pretty, err := prettyJson(raw)
	if err != nil {
		return fmt.Errorf("(%d) format json: %w", year, err)
	}

	outPath := filepath.Join(cfg.RawOutputDir, fmt.Sprintf("jira_%d.json", year))
	if err := storage.Save(pretty, outPath); err != nil {
		return fmt.Errorf("(%d) save file: %w", year, err)
	}

	return nil
}

func paginate(ctx context.Context, client *http.Client, cfg config.Config, reqBuilder func(startAt int) jira.SearchRequest) ([]jira.RawIssue, error) {
	var (
		startAt   = 0
		collected []jira.RawIssue
	)

	for {
		reqBody := reqBuilder(startAt)
		req, err := buildHTTPRequest(ctx, cfg.JiraBaseURL, cfg.JiraApiToken, reqBody)

		if err != nil {
			return nil, fmt.Errorf("paginate build request failed: %w", err)
		}

		_, parsed, err := doRequest(client, req)

		if err != nil {
			return nil, fmt.Errorf("paginate request failed: %w", err)
		}

		collected = append(collected, parsed.Issues...)
		startAt += len(parsed.Issues)

		if startAt >= parsed.Total || len(parsed.Issues) == 0 {
			break
		}

	}

	return collected, nil

}
