package task

import "time"

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type RecurrenceType string

const (
	RecurrenceDaily         RecurrenceType = "daily"
	RecurrenceMonthly       RecurrenceType = "monthly"
	RecurrenceSpecificDates RecurrenceType = "specific_dates"
	RecurrenceEvenDays      RecurrenceType = "even_days"
	RecurrenceOddDays       RecurrenceType = "odd_days"
)

type Task struct {
	ID                      int64           `json:"id"`
	Title                   string          `json:"title"`
	Description             string          `json:"description"`
	Status                  Status          `json:"status"`
	ScheduledAt             *time.Time      `json:"scheduled_at"`
	ParentTaskID            *int64          `json:"parent_task_id"`
	RecurrenceType          *RecurrenceType `json:"recurrence_type"`
	RecurrenceInterval      *int            `json:"recurrence_interval"`
	RecurrenceMonthDays     []int           `json:"recurrence_month_days"`
	RecurrenceSpecificDates []time.Time     `json:"recurrence_specific_dates"`
	NextGenerateDate        *time.Time      `json:"next_generate_date"`
	CreatedAt               time.Time       `json:"created_at"`
	UpdatedAt               time.Time       `json:"updated_at"`
}

func (t *Task) CalculateNextGenerateDate(after time.Time) *time.Time {
	if t.RecurrenceType == nil {
		return nil
	}

	after = after.Truncate(24 * time.Hour)
	var next time.Time

	switch *t.RecurrenceType {
	case RecurrenceDaily:
		interval := 1
		if t.RecurrenceInterval != nil && *t.RecurrenceInterval > 0 {
			interval = *t.RecurrenceInterval
		}
		next = after.AddDate(0, 0, interval)
	case RecurrenceMonthly:
		if len(t.RecurrenceMonthDays) == 0 {
			return nil
		}
		
		// Find the next day in the current month or next months
		for i := 1; i <= 60; i++ {
			candidate := after.AddDate(0, 0, i)
			day := candidate.Day()
			for _, d := range t.RecurrenceMonthDays {
				if d == day {
					return &candidate
				}
			}
		}
		return nil
	case RecurrenceSpecificDates:
		for _, d := range t.RecurrenceSpecificDates {
			d = d.Truncate(24 * time.Hour)
			if d.After(after) {
				if next.IsZero() || d.Before(next) {
					next = d
				}
			}
		}
		if next.IsZero() {
			return nil
		}
	case RecurrenceEvenDays:
		for i := 1; i <= 4; i++ {
			candidate := after.AddDate(0, 0, i)
			if candidate.Day()%2 == 0 {
				return &candidate
			}
		}
	case RecurrenceOddDays:
		for i := 1; i <= 4; i++ {
			candidate := after.AddDate(0, 0, i)
			if candidate.Day()%2 != 0 {
				return &candidate
			}
		}
	}

	return &next
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

func (r RecurrenceType) Valid() bool {
	switch r {
	case RecurrenceDaily, RecurrenceMonthly, RecurrenceSpecificDates, RecurrenceEvenDays, RecurrenceOddDays:
		return true
	default:
		return false
	}
}
