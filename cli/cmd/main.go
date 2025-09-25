package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/leejss/simple-json-server/cli/config"
	"github.com/leejss/simple-json-server/cli/jira"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	ctx := context.Background()
	client := &http.Client{}
	builder := &jira.JQLQueryBuilder{}
	years := []int{2023, 2024, 2025}

	for _, year := range years {
		if err := processYear(ctx, client, *cfg, builder, year); err != nil {
			fmt.Printf("[%d] 처리 실패: %v\n", year, err)
			continue
		}
		fmt.Println("처리 완료")
	}
}
