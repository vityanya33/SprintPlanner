package jira

import (
	"fmt"

	jiraapi "github.com/andygrunwald/go-jira"
)

// ExportOptions содержит настройки для экспорта
type ExportOptions struct {
	JQL             string `json:"jql"`
	OutputFile      string `json:"output_file"`
	MaxResults      int    `json:"max_results"`
	IncludeComments bool   `json:"include_comments"`
	StartAt         int    `json:"start_at"`
}

// SearchIssuesWithPagination выполняет поиск issues с пагинацией
func (c *Client) SearchIssuesWithPagination(jql string, startAt, maxResults int) ([]*jiraapi.Issue, int, error) {
	searchOptions := &jiraapi.SearchOptions{
		StartAt:    startAt,
		MaxResults: maxResults,
		Fields:     []string{"summary", "description", "timeTracking", "duedate"},
		Expand:     "timeTracking",
	}

	issues, resp, err := c.api.Issue.Search(jql, searchOptions)
	if err != nil {
		return nil, 0, NewGetIssueError(err)
	}

	result := make([]*jiraapi.Issue, len(issues))
	for i := range issues {
		result[i] = &issues[i]
	}

	return result, resp.Total, nil
}

// GetIssueWithFullDetails получает issue со всеми деталями
func (c *Client) GetIssueWithFullDetails(issueKey string) (*jiraapi.Issue, error) {
	// Определяем expand параметры
	expandParams := "names,schema,operations,editmeta,changelog,renderedFields"

	// Получаем issue с расширенными полями
	issue, _, err := c.api.Issue.Get(issueKey, &jiraapi.GetQueryOptions{
		Expand: expandParams,
	})
	if err != nil {
		return nil, NewGetIssueError(err)
	}

	return issue, nil
}

// ExportIssues экспортирует issues по JQL запросу в JSON файл
func (c *Client) ExportIssues(options ExportOptions) ([]*jiraapi.Issue, error) {
	fmt.Printf("Starting export with JQL: %s", options.JQL)

	// Устанавливаем значения по умолчанию
	if options.MaxResults <= 0 {
		options.MaxResults = 50 // Jira default
	}

	var allExportedIssues []*jiraapi.Issue
	startAt := options.StartAt

	for {
		fmt.Printf("Fetching issues: startAt=%d, maxResults=%d", startAt, options.MaxResults)

		// Получаем batch issues
		issues, total, err := c.SearchIssuesWithPagination(options.JQL, startAt, options.MaxResults)
		if err != nil {
			return nil, NewGetIssueError(err)
		}

		if len(issues) == 0 {
			break
		}

		// Обрабатываем каждый issue
		for _, issue := range issues {
			exportedIssue, err := c.GetIssueWithFullDetails(issue.Key)
			if err != nil {
				return nil, NewGetIssueError(err)
			}

			allExportedIssues = append(allExportedIssues, exportedIssue)

			fmt.Printf("Exported issue: %s (%d/%d)", issue.Key, len(allExportedIssues), total)
		}

		// Проверяем, есть ли еще issues
		startAt += len(issues)
		if startAt >= total {
			break
		}
	}

	return allExportedIssues, nil
}
