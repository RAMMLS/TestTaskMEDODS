package task

import (
	"context"
	"log"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

func (s *Service) GenerateRecurringTasks(ctx context.Context) error {
	now := s.now()
	tasks, err := s.repo.GetDueRecurringTasks(ctx, now)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if err := s.generateNextTask(ctx, &task); err != nil {
			log.Printf("Failed to generate task for template %d: %v", task.ID, err)
			continue
		}
	}

	return nil
}

func (s *Service) generateNextTask(ctx context.Context, template *taskdomain.Task) error {
	// Calculate the next generation date first
	nextDate := template.NextGenerateDate
	if nextDate == nil {
		return nil
	}

	// Create a new task instance
	newTask := &taskdomain.Task{
		Title:       template.Title,
		Description: template.Description,
		Status:      taskdomain.StatusNew,
		ScheduledAt: nextDate,
		ParentTaskID: &template.ID,
		CreatedAt:   s.now(),
		UpdatedAt:   s.now(),
	}

	_, err := s.repo.Create(ctx, newTask)
	if err != nil {
		return err
	}

	// Update the template's NextGenerateDate
	template.NextGenerateDate = template.CalculateNextGenerateDate(*nextDate)
	template.UpdatedAt = s.now()
	_, err = s.repo.Update(ctx, template)
	if err != nil {
		return err
	}

	return nil
}
