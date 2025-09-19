package jira

import "fmt"

type JQLQueryBuilder struct{}

func (q *JQLQueryBuilder) SearchByYear(year int, assignee string) string {
	return fmt.Sprintf("assignee = %s AND created >= %d-01-01 AND created < %d-01-01 order by created ASC", assignee, year, year+1)
}
