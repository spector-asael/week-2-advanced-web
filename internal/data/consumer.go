package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Consumer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	Version   int32     `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type ConsumerModel struct {
	DB *sql.DB
}

func (m *ConsumerModel) Insert(consumer *Consumer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = tx.QueryRowContext(ctx, `
		INSERT INTO consumers (name, email, status)
		VALUES ($1, $2, $3)
		RETURNING id, version, created_at, updated_at`,
		consumer.Name,
		consumer.Email,
		consumer.Status,
	).Scan(
		&consumer.ID,
		&consumer.Version,
		&consumer.CreatedAt,
		&consumer.UpdatedAt,
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

func (m *ConsumerModel) GetById(id string) (*Consumer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var consumer Consumer

	err := m.DB.QueryRowContext(ctx, `
		SELECT id, name, email, status, version, updated_at, created_at
		FROM consumers
		WHERE id = $1`,
		id,
	).Scan(
		&consumer.ID,
		&consumer.Name,
		&consumer.Email,
		&consumer.Status,
		&consumer.Version,
		&consumer.UpdatedAt,
		&consumer.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}

	return &consumer, nil
}

func (m *ConsumerModel) GetAll() ([]*Consumer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, `
		SELECT id, name, email, status, version, updated_at, created_at
		FROM consumers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consumers []*Consumer

	for rows.Next() {
		var consumer Consumer

		err := rows.Scan(
			&consumer.ID,
			&consumer.Name,
			&consumer.Email,
			&consumer.Status,
			&consumer.Version,
			&consumer.UpdatedAt,
			&consumer.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		consumers = append(consumers, &consumer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return consumers, nil
}

func (m *ConsumerModel) UpdateById(consumer *Consumer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE consumers
		SET name = $1,
		    email = $2,
		    status = $3,
		    version = version + 1,
		    updated_at = NOW()
		WHERE id = $4`,
		consumer.Name,
		consumer.Email,
		consumer.Status,
		consumer.ID,
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
