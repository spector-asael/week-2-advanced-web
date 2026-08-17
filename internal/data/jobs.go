package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Jobs struct {
	ID           string     `json:"id"`
	ConsumerID   string     `json:"consumer_id"`
	JobType      string     `json:"job_type"`
	Status       string     `json:"status"`
	Payload      string     `json:"payload"`
	Result       *string    `json:"result"`
	ErrorMessage *string    `json:"error_message"`
	StartedAt    *time.Time `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

type JobsModel struct {
	DB *sql.DB
}

func (m *JobsModel) Insert(job *Jobs) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = tx.QueryRowContext(ctx, `
		INSERT INTO jobs (
			consumer_id,
			job_type,
			status,
			payload,
			result,
			started_at,
			completed_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`,
		job.ConsumerID,
		job.JobType,
		job.Status,
		job.Payload,
		job.Result,
		job.StartedAt,
		job.CompletedAt,
	).Scan(
		&job.ID,
		&job.CreatedAt,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *JobsModel) GetById(id string) (*Jobs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var job Jobs

	err := m.DB.QueryRowContext(ctx, `
		SELECT id, consumer_id, job_type, status, payload, result,
		       error_message, started_at, completed_at, created_at
		FROM jobs
		WHERE id = $1`,
		id,
	).Scan(
		&job.ID,
		&job.ConsumerID,
		&job.JobType,
		&job.Status,
		&job.Payload,
		&job.Result,
		&job.ErrorMessage,
		&job.StartedAt,
		&job.CompletedAt,
		&job.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("job not found")
		}
		return nil, err
	}

	return &job, nil
}

func (m *JobsModel) UpdateById(job *Jobs) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE jobs
		SET consumer_id = $1,
		    job_type = $2,
		    status = $3,
		    payload = $4,
		    result = $5,
		    error_message = $6,
		    started_at = $7,
		    completed_at = $8
		WHERE id = $9`,
		job.ConsumerID,
		job.JobType,
		job.Status,
		job.Payload,
		job.Result,
		job.ErrorMessage,
		job.StartedAt,
		job.CompletedAt,
		job.ID,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *JobsModel) GetAll() ([]*Jobs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, `
		SELECT id, consumer_id, job_type, status, payload, result,
		       error_message, started_at, completed_at, created_at
		FROM jobs`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*Jobs

	for rows.Next() {
		var job Jobs

		err := rows.Scan(
			&job.ID,
			&job.ConsumerID,
			&job.JobType,
			&job.Status,
			&job.Payload,
			&job.Result,
			&job.ErrorMessage,
			&job.StartedAt,
			&job.CompletedAt,
			&job.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, &job)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}
