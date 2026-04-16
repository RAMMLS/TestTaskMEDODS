package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `
		INSERT INTO tasks (
			title, description, status, scheduled_at, parent_task_id, 
			recurrence_type, recurrence_interval, recurrence_month_days, recurrence_specific_dates, next_generate_date,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, title, description, status, scheduled_at, parent_task_id, recurrence_type, recurrence_interval, recurrence_month_days, recurrence_specific_dates, next_generate_date, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query, 
		task.Title, task.Description, task.Status, task.ScheduledAt, task.ParentTaskID,
		task.RecurrenceType, task.RecurrenceInterval, task.RecurrenceMonthDays, task.RecurrenceSpecificDates, task.NextGenerateDate,
		task.CreatedAt, task.UpdatedAt,
	)
	created, err := scanTask(row)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	const query = `
		SELECT id, title, description, status, scheduled_at, parent_task_id, recurrence_type, recurrence_interval, recurrence_month_days, recurrence_specific_dates, next_generate_date, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	found, err := scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}

		return nil, err
	}

	return found, nil
}

func (r *Repository) Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `
		UPDATE tasks
		SET title = $1,
			description = $2,
			status = $3,
			scheduled_at = $4,
			parent_task_id = $5,
			recurrence_type = $6,
			recurrence_interval = $7,
			recurrence_month_days = $8,
			recurrence_specific_dates = $9,
			next_generate_date = $10,
			updated_at = $11
		WHERE id = $12
		RETURNING id, title, description, status, scheduled_at, parent_task_id, recurrence_type, recurrence_interval, recurrence_month_days, recurrence_specific_dates, next_generate_date, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query, 
		task.Title, task.Description, task.Status, task.ScheduledAt, task.ParentTaskID,
		task.RecurrenceType, task.RecurrenceInterval, task.RecurrenceMonthDays, task.RecurrenceSpecificDates, task.NextGenerateDate,
		task.UpdatedAt, task.ID,
	)
	updated, err := scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}

		return nil, err
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM tasks WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return taskdomain.ErrNotFound
	}

	return nil
}

func (r *Repository) List(ctx context.Context) ([]taskdomain.Task, error) {
	const query = `
		SELECT id, title, description, status, scheduled_at, parent_task_id, recurrence_type, recurrence_interval, recurrence_month_days, recurrence_specific_dates, next_generate_date, created_at, updated_at
		FROM tasks
		ORDER BY id DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]taskdomain.Task, 0)
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, *task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) GetDueRecurringTasks(ctx context.Context, before time.Time) ([]taskdomain.Task, error) {
	const query = `
		SELECT id, title, description, status, scheduled_at, parent_task_id, recurrence_type, recurrence_interval, recurrence_month_days, recurrence_specific_dates, next_generate_date, created_at, updated_at
		FROM tasks
		WHERE recurrence_type IS NOT NULL 
		  AND next_generate_date IS NOT NULL 
		  AND next_generate_date <= $1
	`

	rows, err := r.pool.Query(ctx, query, before)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]taskdomain.Task, 0)
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, *task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

type taskScanner interface {
	Scan(dest ...any) error
}

func scanTask(scanner taskScanner) (*taskdomain.Task, error) {
	var (
		task   taskdomain.Task
		status string
		rtype  *string
	)

	if err := scanner.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&status,
		&task.ScheduledAt,
		&task.ParentTaskID,
		&rtype,
		&task.RecurrenceInterval,
		&task.RecurrenceMonthDays,
		&task.RecurrenceSpecificDates,
		&task.NextGenerateDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		return nil, err
	}

	task.Status = taskdomain.Status(status)
	if rtype != nil {
		rt := taskdomain.RecurrenceType(*rtype)
		task.RecurrenceType = &rt
	}

	return &task, nil
}
