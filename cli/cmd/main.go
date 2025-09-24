package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/leejss/simple-json-server/cli/config"
	"github.com/leejss/simple-json-server/cli/internal/storage"
	"github.com/leejss/simple-json-server/cli/jira"
)

func main() {
	config, err := config.LoadConfig()

	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	queryBuilder := &jira.JQLQueryBuilder{}

	years := []int{2023, 2024, 2025}

	// HTTP 클라이언트를 한 번만 생성하여 재사용
	client := &http.Client{}

	// jiraClient := jira.NewJiraClient(config.JiraBaseURL, config.JiraApiToken)

	for _, year := range years {
		// 각 연도별 JQL 생성
		jqlQuery := queryBuilder.SearchByYear(year, config.Username)

		// 요청 페이로드 구성
		reqBody := jira.SearchRequest{
			JQL:        jqlQuery,
			StartAt:    0,
			MaxResults: 100,
			Fields:     []string{"key", "summary", "created", "description"},
		}

		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", config.JiraBaseURL+"/rest/api/2/search", bytes.NewBuffer(jsonBody))

		// 인증 및 헤더 설정
		req.Header.Set("Authorization", "Bearer "+config.JiraApiToken)
		req.Header.Set("Content-Type", "application/json")

		// 요청 실행
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("[%d] 요청 오류: %v\n", year, err)
			continue // 다음 연도로 계속
		}

		// 응답 바디는 반드시 close
		func() {
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("[%d] 상태 오류: %s\n", year, string(body))
				return
			}

			// JSON 예쁘게 포맷팅
			var prettyJson bytes.Buffer
			if err := json.Indent(&prettyJson, body, "", "  "); err != nil {
				fmt.Printf("[%d] JSON 포맷 오류: %v\n", year, err)
				return
			}

			// 연도별 파일 경로 생성 후 저장
			outPath := filepath.Join(config.RawOutputDir, fmt.Sprintf("jira_%d.json", year))
			if err := storage.Save(prettyJson.Bytes(), outPath); err != nil {
				fmt.Printf("[%d] 저장 오류: %v\n", year, err)
				return
			}
		}()
	}

}
