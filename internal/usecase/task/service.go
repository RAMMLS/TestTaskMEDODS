package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		Title:                   normalized.Title,
		Description:             normalized.Description,
		Status:                  normalized.Status,
		ScheduledAt:             normalized.ScheduledAt,
		RecurrenceType:          normalized.RecurrenceType,
		RecurrenceInterval:      normalized.RecurrenceInterval,
		RecurrenceMonthDays:     normalized.RecurrenceMonthDays,
		RecurrenceSpecificDates: normalized.RecurrenceSpecificDates,
	}
	now := s.now()
	model.CreatedAt = now
	model.UpdatedAt = now

	if model.RecurrenceType != nil {
		if model.ScheduledAt == nil {
			model.ScheduledAt = &now
		}
		model.NextGenerateDate = model.CalculateNextGenerateDate(*model.ScheduledAt)
	}

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	task.Title = normalized.Title
	task.Description = normalized.Description
	task.Status = normalized.Status
	task.ScheduledAt = normalized.ScheduledAt
	task.RecurrenceType = normalized.RecurrenceType
	task.RecurrenceInterval = normalized.RecurrenceInterval
	task.RecurrenceMonthDays = normalized.RecurrenceMonthDays
	task.RecurrenceSpecificDates = normalized.RecurrenceSpecificDates
	task.UpdatedAt = s.now()

	if task.RecurrenceType != nil {
		if task.ScheduledAt == nil {
			now := s.now()
			task.ScheduledAt = &now
		}
		task.NextGenerateDate = task.CalculateNextGenerateDate(*task.ScheduledAt)
	} else {
		task.NextGenerateDate = nil
	}

	updated, err := s.repo.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if input.RecurrenceType != nil {
		if !input.RecurrenceType.Valid() {
			return CreateInput{}, fmt.Errorf("%w: invalid recurrence type", ErrInvalidInput)
		}
		if *input.RecurrenceType == taskdomain.RecurrenceMonthly && len(input.RecurrenceMonthDays) == 0 {
			return CreateInput{}, fmt.Errorf("%w: recurrence_month_days is required for monthly recurrence", ErrInvalidInput)
		}
		if *input.RecurrenceType == taskdomain.RecurrenceSpecificDates && len(input.RecurrenceSpecificDates) == 0 {
			return CreateInput{}, fmt.Errorf("%w: recurrence_specific_dates is required for specific_dates recurrence", ErrInvalidInput)
		}
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if input.RecurrenceType != nil {
		if !input.RecurrenceType.Valid() {
			return UpdateInput{}, fmt.Errorf("%w: invalid recurrence type", ErrInvalidInput)
		}
		if *input.RecurrenceType == taskdomain.RecurrenceMonthly && len(input.RecurrenceMonthDays) == 0 {
			return UpdateInput{}, fmt.Errorf("%w: recurrence_month_days is required for monthly recurrence", ErrInvalidInput)
		}
		if *input.RecurrenceType == taskdomain.RecurrenceSpecificDates && len(input.RecurrenceSpecificDates) == 0 {
			return UpdateInput{}, fmt.Errorf("%w: recurrence_specific_dates is required for specific_dates recurrence", ErrInvalidInput)
		}
	}

	return input, nil
}
