package models

type RawIssue struct {
	Key    string `json:"key"`
	Fields struct {
		Summary     string `json:"summary"`
		Created     string `json:"created"`
		Description string `json:"description"`
	} `json:"fields"`
}

// Issue는 포맷팅된 간소화 버전
type Issue struct {
	Summary     string `json:"summary"`
	Created     string `json:"created"`
	Description string `json:"description"`
}

// Define response structures

// IssuesResponse는 API 응답 구조
type IssuesResponse struct {
	Year   int     `json:"year,omitempty"`
	Total  int     `json:"total"`
	Issues []Issue `json:"issues"`
}

// ErrorResponse는 에러 응답 구조
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
