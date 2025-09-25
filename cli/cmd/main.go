package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/leejss/simple-json-server/cli/internal/storage"
	"github.com/leejss/simple-json-server/cli/jira"
)

func main() {
	// cfg, err := config.LoadConfig()

	// if err != nil {
	// 	fmt.Println("Error loading .env file")
	// 	return
	// }

	// ctx := context.Background()
	// client := &http.Client{}
	// builder := &jira.JQLQueryBuilder{}
	// years := []int{2023, 2024, 2025}

	// for _, year := range years {
	// 	if err := processYear(ctx, client, *cfg, builder, year); err != nil {
	// 		fmt.Printf("[%d] 처리 실패: %v\n", year, err)
	// 		continue
	// 	}
	// 	fmt.Println("처리 완료")
	// }

	// READ output/raw dir -> iterate file validate 이게 정말 내가 생성한 json인가 -> 틀리다면 skip. 맞다면 processing 진행 -> json을 읽고 구조체로 변환 -> 이 구조체에서 특정 필드를 선택 후 formatted struct 생성 -> 생성한 구조체를 다시 json으로 직렬화 -> output/formatted dir에 저장

	// listRawFiles
	files, err := listRawFiles("../output/raw")
	if err != nil {
		fmt.Printf("Error listing raw files: %v\n", err)
		return
	}

	var validFiles []string

	for _, file := range files {
		fileName := filepath.Base(file)
		if !isGeneratedFile(fileName) {
			fmt.Println("Not a generated file", fileName)
			continue
		}

		validFiles = append(validFiles, file)
	}

	fmt.Println("validFiles", validFiles)

	for _, file := range validFiles {
		parsed, err := parseRawJson(file)
		if err != nil {
			fmt.Printf("Error parsing raw json: %v\n", err)
			continue
		}

		formatted, err := buildFormattedIssue(parsed)

		if err != nil {
			fmt.Printf("Error building formatted issue: %v\n", err)
			continue
		}

		// create filepath
		fileName := filepath.Base(file)
		filePath := filepath.Join("../output/formatted", strings.TrimPrefix(fileName, "jira_"))

		if err := saveFormattedIssues(formatted, filePath); err != nil {
			fmt.Printf("Error saving formatted issue: %v\n", err)
			continue
		}
	}

}

type Issue struct {
	Summary     string `json:"summary"`
	Created     string `json:"created"`
	Description string `json:"description"`
}

func listRawFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read raw dir: %w", err)
	}

	var paths []string
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		paths = append(paths, filepath.Join(dir, entry.Name()))
	}

	return paths, nil
}

func isGeneratedFile(fileName string) bool {
	return strings.HasPrefix(fileName, "jira_")
}

func parseRawJson(filePath string) ([]jira.RawIssue, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read raw json: %w", err)
	}
	var rawIssues []jira.RawIssue
	if err := json.Unmarshal(data, &rawIssues); err != nil {
		return nil, fmt.Errorf("failed to unmarshal raw json: %w", err)
	}
	return rawIssues, nil
}

func buildFormattedIssue(rawIssue []jira.RawIssue) ([]Issue, error) {
	var issues []Issue
	for _, issue := range rawIssue {
		issues = append(issues, Issue{
			Summary:     issue.Fields.Summary,
			Created:     issue.Fields.Created,
			Description: issue.Fields.Description,
		})
	}

	return issues, nil
}

func saveFormattedIssues(issues []Issue, dir string) error {

	data, err := json.Marshal(issues)
	if err != nil {
		return fmt.Errorf("failed to marshal issues: %w", err)
	}

	return storage.Save(data, dir)
}
