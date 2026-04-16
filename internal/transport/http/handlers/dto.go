package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type taskMutationDTO struct {
	Title                   string                     `json:"title"`
	Description             string                     `json:"description"`
	Status                  taskdomain.Status          `json:"status"`
	ScheduledAt             *time.Time                 `json:"scheduled_at"`
	RecurrenceType          *taskdomain.RecurrenceType `json:"recurrence_type"`
	RecurrenceInterval      *int                       `json:"recurrence_interval"`
	RecurrenceMonthDays     []int                      `json:"recurrence_month_days"`
	RecurrenceSpecificDates []time.Time                `json:"recurrence_specific_dates"`
}

type taskDTO struct {
	ID                      int64                      `json:"id"`
	Title                   string                     `json:"title"`
	Description             string                     `json:"description"`
	Status                  taskdomain.Status          `json:"status"`
	ScheduledAt             *time.Time                 `json:"scheduled_at"`
	ParentTaskID            *int64                     `json:"parent_task_id"`
	RecurrenceType          *taskdomain.RecurrenceType `json:"recurrence_type"`
	RecurrenceInterval      *int                       `json:"recurrence_interval"`
	RecurrenceMonthDays     []int                      `json:"recurrence_month_days"`
	RecurrenceSpecificDates []time.Time                `json:"recurrence_specific_dates"`
	NextGenerateDate        *time.Time                 `json:"next_generate_date"`
	CreatedAt               time.Time                  `json:"created_at"`
	UpdatedAt               time.Time                  `json:"updated_at"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:                      task.ID,
		Title:                   task.Title,
		Description:             task.Description,
		Status:                  task.Status,
		ScheduledAt:             task.ScheduledAt,
		ParentTaskID:            task.ParentTaskID,
		RecurrenceType:          task.RecurrenceType,
		RecurrenceInterval:      task.RecurrenceInterval,
		RecurrenceMonthDays:     task.RecurrenceMonthDays,
		RecurrenceSpecificDates: task.RecurrenceSpecificDates,
		NextGenerateDate:        task.NextGenerateDate,
		CreatedAt:               task.CreatedAt,
		UpdatedAt:               task.UpdatedAt,
	}
}
