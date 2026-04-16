CREATE TABLE IF NOT EXISTS tasks (
	id BIGSERIAL PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT NOT NULL DEFAULT '',
	status TEXT NOT NULL,
	scheduled_at TIMESTAMPTZ,
	parent_task_id BIGINT REFERENCES tasks(id) ON DELETE SET NULL,
	recurrence_type TEXT,
	recurrence_interval INT,
	recurrence_month_days INT[],
	recurrence_specific_dates DATE[],
	next_generate_date DATE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks (status);
CREATE INDEX IF NOT EXISTS idx_tasks_next_generate_date ON tasks (next_generate_date);
