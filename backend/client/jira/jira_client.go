package jira

import (
	jiraapi "github.com/andygrunwald/go-jira"
)

// Client wraps the go-jira client for easier usage and extension.
type Client struct {
	api *jiraapi.Client
}

// NewClient создает нового клиента Jira.
func NewClient(baseURL, apiToken string) (*Client, error) {
	tp := jiraapi.PATAuthTransport{
		Token: apiToken,
	}
	apiClient, err := jiraapi.NewClient(tp.Client(), baseURL)
	if err != nil {
		return nil, NewCreateClientError(err)
	}
	return &Client{api: apiClient}, nil
}

// NewJiraClientFromEnv создает клиента Jira, используя переменные окружения:
//
//	JIRA_BASE_URL, JIRA_USERNAME, JIRA_API_TOKEN
func NewJiraClient(baseURL, apiToken string) (*Client, error) {
	return NewClient(baseURL, apiToken)
}

// GetIssue получает тикет по ключу.
func (c *Client) GetIssueSummary(issueKey string) (*jiraapi.Issue, error) {
	issue, resp, err := c.api.Issue.Get(issueKey, nil)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return nil, NewNotFoundError()
		}
		if err.Error() == "404" || err.Error() == "jira: 404" {
			return nil, NewNotFoundError()
		}
		return nil, NewGetIssueError(err)
	}
	return issue, nil
}

// AddComment posts a comment to the specified Jira issue.
func (c *Client) AddComment(issueKey, comment string) error {
	_, _, err := c.api.Issue.AddComment(issueKey, &jiraapi.Comment{Body: comment})
	if err != nil {
		return NewAddCommentError(err)
	}
	return nil
}

// SearchIssues searches for issues using the provided JQL query.
func (c *Client) SearchIssues(jql string) ([]*jiraapi.Issue, error) {
	issues, _, err := c.api.Issue.Search(jql, nil)
	if err != nil {
		return nil, NewGetIssueError(err)
	}
	result := make([]*jiraapi.Issue, len(issues))
	for i := range issues {
		result[i] = &issues[i]
	}
	return result, nil
}
