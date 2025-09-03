package services

import (
	"backend/client/jira"
	"backend/models"
	"backend/repositories"
	"context"
	"fmt"
	"time"
)

type TaskService interface {
	GetAllTasks(ctx context.Context) ([]models.Task, error)
	GetTaskByID(ctx context.Context, id string) (*models.Task, error)
	CreateTask(ctx context.Context, task *models.Task) error
	UpdateTask(ctx context.Context, id string, task *models.Task) error
	DeleteTask(ctx context.Context, id string) error
	UpdateTaskUsers(ctx context.Context, taskID string, userIDs []int) error
	GetAvailableUsers(ctx context.Context, hours int) ([]models.UserWithLoad, error)
	SyncTasksWithJira(ctx context.Context, jql string) error
}

// Исправьте на TaskServiceImpl вместо taskService
type TaskServiceImpl struct {
	jiraClient *jira.Client
	taskRepo   repositories.TaskRepository
}

// И возвращайте интерфейс TaskService
func NewTaskService(taskRepo repositories.TaskRepository, jiraClient *jira.Client) TaskService {
	return &TaskServiceImpl{taskRepo: taskRepo, jiraClient: jiraClient}
}

func (s *TaskServiceImpl) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	return s.taskRepo.GetAllTasks(ctx)
}

func (s *TaskServiceImpl) GetTaskByID(ctx context.Context, id string) (*models.Task, error) {
	return s.taskRepo.GetTaskByID(ctx, id)
}

func (s *TaskServiceImpl) CreateTask(ctx context.Context, task *models.Task) error {
	return s.taskRepo.CreateTask(ctx, task)
}

func (s *TaskServiceImpl) UpdateTask(ctx context.Context, id string, task *models.Task) error {
	return s.taskRepo.UpdateTask(ctx, id, task)
}

func (s *TaskServiceImpl) DeleteTask(ctx context.Context, id string) error {
	return s.taskRepo.DeleteTask(ctx, id)
}

func (s *TaskServiceImpl) UpdateTaskUsers(ctx context.Context, taskID string, userIDs []int) error {
	return s.taskRepo.UpdateTaskUsers(ctx, taskID, userIDs)
}

func (s *TaskServiceImpl) GetAvailableUsers(ctx context.Context, hours int) ([]models.UserWithLoad, error) {
	users, err := s.taskRepo.GetAvailableUsers(ctx)
	if err != nil {
		return nil, err
	}

	var available []models.UserWithLoad
	for _, user := range users {
		// Простой расчет: свободный ресурс = общий ресурс - занятость
		// Фильтруем тех, у кого свободного ресурса >= требуемого ресурса задачи
		if user.Free >= hours {
			available = append(available, user)
		}
	}

	return available, nil
}

func (s *TaskServiceImpl) SyncTasksWithJira(ctx context.Context, jql string) error {
	if s.jiraClient == nil {
		return fmt.Errorf("jira client is not configured")
	}

	const pageSize = 50

	issues, err := s.jiraClient.ExportIssues(jira.ExportOptions{JQL: jql, StartAt: 0, MaxResults: pageSize})
	if err != nil {
		return err
	}
	if len(issues) == 0 {
		return nil
	}

	batch := make([]models.Task, 0, len(issues))
	for _, issue := range issues {
		fmt.Printf("Importing Jira issue: key=%s, summary=%s", issue.Key, issue.Fields.Summary)
		title := ""
		if issue.Fields != nil && issue.Fields.Summary != "" {
			title = issue.Fields.Summary
		}

		hoursDuration := 0
		if issue.Fields != nil && issue.Fields.TimeTracking != nil {
			if issue.Fields.TimeTracking.OriginalEstimateSeconds > 0 {
				hoursDuration = issue.Fields.TimeTracking.OriginalEstimateSeconds / 3600
			}
		}

		deadline := time.Time{}
		if issue.Fields != nil {
			if !time.Time(issue.Fields.Duedate).IsZero() {
				deadline = time.Time(issue.Fields.Duedate)
			} else {
				deadline = time.Now()
			}
		}

		startDate := calculateStartDateByWorkHours(deadline, hoursDuration)
		task := models.Task{
			ID:        issue.Key,
			Title:     title,
			Hours:     hoursDuration,
			StartDate: startDate,
			Deadline:  deadline,
		}
		batch = append(batch, task)
	}

	if err := s.taskRepo.BulkCreateOrUpdateTasks(ctx, batch); err != nil {
		return err
	}
	return nil
}

// calculateStartDateByWorkHours returns the start date by subtracting working hours from the due date.
// Assumptions:
// - Working day has 8 hours.
// - Weekends (Saturday, Sunday) are non-working and are skipped.
// - Dates are stored as DATE (no time component). If due date falls on a weekend, it is adjusted back to the previous Friday.
func calculateStartDateByWorkHours(dueDate time.Time, hours int) time.Time {
	if dueDate.IsZero() || hours <= 0 {
		return dueDate
	}

	// Adjust due date to previous weekday if it is on weekend
	d := dueDate
	for d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
		d = d.AddDate(0, 0, -1)
	}

	// Compute number of working days needed (ceil by 8 hours per day)
	daysNeeded := (hours + 7) / 8
	// We count the due date as the last working day, so move back (daysNeeded-1) working days
	steps := daysNeeded - 1
	for steps > 0 {
		d = d.AddDate(0, 0, -1)
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			continue
		}
		steps--
	}
	return d
}
