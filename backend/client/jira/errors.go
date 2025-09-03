package jira

import "fmt"

// NotFoundError represents a Jira ticket not found error.
type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return "jira ticket not found (404)"
}

// MissingEnvVarsError represents missing environment variables error.
type MissingEnvVarsError struct{}

func (e *MissingEnvVarsError) Error() string {
	return "environment variables JIRA_BASE_URL, JIRA_USERNAME, and JIRA_API_TOKEN must be set"
}

// CreateClientError wraps errors from client creation.
type CreateClientError struct {
	Err error
}

func (e *CreateClientError) Error() string {
	return fmt.Sprintf("failed to create Jira client: %v", e.Err)
}

func (e *CreateClientError) Unwrap() error { return e.Err }

// GetIssueError wraps errors from getting an issue.
type GetIssueError struct {
	Err error
}

func (e *GetIssueError) Error() string {
	return fmt.Sprintf("failed to get Jira issue: %v", e.Err)
}

func (e *GetIssueError) Unwrap() error { return e.Err }

// AddCommentError wraps errors from adding a comment to an issue.
type AddCommentError struct {
	Err error
}

func (e *AddCommentError) Error() string {
	return fmt.Sprintf("failed to add comment to Jira issue: %v", e.Err)
}

// Constructors
func NewNotFoundError() error              { return &NotFoundError{} }
func NewMissingEnvVarsError() error        { return &MissingEnvVarsError{} }
func NewCreateClientError(err error) error { return &CreateClientError{Err: err} }
func NewGetIssueError(err error) error     { return &GetIssueError{Err: err} }
func NewAddCommentError(err error) error   { return &AddCommentError{Err: err} }

// Error type checkers
func IsNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

func IsMissingEnvVarsError(err error) bool {
	_, ok := err.(*MissingEnvVarsError)
	return ok
}

func IsCreateClientError(err error) bool {
	_, ok := err.(*CreateClientError)
	return ok
}

func IsGetIssueError(err error) bool {
	_, ok := err.(*GetIssueError)
	return ok
}

func IsAddCommentError(err error) bool {
	_, ok := err.(*AddCommentError)
	return ok
}
