package task_test

import (
	"testing"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

func TestCalculateNextGenerateDate_Daily(t *testing.T) {
	rtype := taskdomain.RecurrenceDaily
	interval := 2
	task := taskdomain.Task{
		RecurrenceType:     &rtype,
		RecurrenceInterval: &interval,
	}

	after := time.Date(2026, 3, 24, 0, 0, 0, 0, time.UTC)
	next := task.CalculateNextGenerateDate(after)

	if next == nil {
		t.Fatal("expected next date, got nil")
	}

	expected := time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Fatalf("expected %v, got %v", expected, next)
	}
}

func TestCalculateNextGenerateDate_Monthly(t *testing.T) {
	rtype := taskdomain.RecurrenceMonthly
	task := taskdomain.Task{
		RecurrenceType:      &rtype,
		RecurrenceMonthDays: []int{1, 15},
	}

	after := time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC)
	next := task.CalculateNextGenerateDate(after)

	if next == nil {
		t.Fatal("expected next date, got nil")
	}

	expected := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Fatalf("expected %v, got %v", expected, next)
	}
}

func TestCalculateNextGenerateDate_SpecificDates(t *testing.T) {
	rtype := taskdomain.RecurrenceSpecificDates
	task := taskdomain.Task{
		RecurrenceType: &rtype,
		RecurrenceSpecificDates: []time.Time{
			time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
		},
	}

	after := time.Date(2026, 3, 24, 0, 0, 0, 0, time.UTC)
	next := task.CalculateNextGenerateDate(after)

	if next == nil {
		t.Fatal("expected next date, got nil")
	}

	expected := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Fatalf("expected %v, got %v", expected, next)
	}
}

func TestCalculateNextGenerateDate_EvenDays(t *testing.T) {
	rtype := taskdomain.RecurrenceEvenDays
	task := taskdomain.Task{
		RecurrenceType: &rtype,
	}

	after := time.Date(2026, 3, 24, 0, 0, 0, 0, time.UTC) // 24th
	next := task.CalculateNextGenerateDate(after)

	if next == nil {
		t.Fatal("expected next date, got nil")
	}

	expected := time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC) // 26th
	if !next.Equal(expected) {
		t.Fatalf("expected %v, got %v", expected, next)
	}
}

func TestCalculateNextGenerateDate_OddDays(t *testing.T) {
	rtype := taskdomain.RecurrenceOddDays
	task := taskdomain.Task{
		RecurrenceType: &rtype,
	}

	after := time.Date(2026, 3, 24, 0, 0, 0, 0, time.UTC) // 24th
	next := task.CalculateNextGenerateDate(after)

	if next == nil {
		t.Fatal("expected next date, got nil")
	}

	expected := time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC) // 25th
	if !next.Equal(expected) {
		t.Fatalf("expected %v, got %v", expected, next)
	}
}
