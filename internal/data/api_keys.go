package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type API_Keys struct {
	ID         string     `json:"id"`
	ConsumerID string     `json:"consumer_id"`
	KeyHash    string     `json:"key_hash"`
	KeyPrefix  string     `json:"key_prefix"`
	Status     string     `json:"status"`
	LastUsedAt *time.Time `json:"last_used_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

type API_KeysModel struct {
	DB *sql.DB
}

func (m *API_KeysModel) Insert(apiKey *API_Keys) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = tx.QueryRowContext(ctx, `
		INSERT INTO api_keys (
			consumer_id,
			key_hash,
			key_prefix,
			status,
			last_used_at,
			expires_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`,
		apiKey.ConsumerID,
		apiKey.KeyHash,
		apiKey.KeyPrefix,
		apiKey.Status,
		apiKey.LastUsedAt,
		apiKey.ExpiresAt,
	).Scan(
		&apiKey.ID,
		&apiKey.CreatedAt,
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

func (m *API_KeysModel) GetById(id string) (*API_Keys, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var apiKey API_Keys

	err := m.DB.QueryRowContext(ctx, `
		SELECT id, consumer_id, key_hash, key_prefix, status,
		       last_used_at, expires_at, created_at
		FROM api_keys
		WHERE id = $1`,
		id,
	).Scan(
		&apiKey.ID,
		&apiKey.ConsumerID,
		&apiKey.KeyHash,
		&apiKey.KeyPrefix,
		&apiKey.Status,
		&apiKey.LastUsedAt,
		&apiKey.ExpiresAt,
		&apiKey.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}

	return &apiKey, nil
}

func (m *API_KeysModel) GetAll() ([]*API_Keys, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, `
		SELECT id, consumer_id, key_hash, key_prefix, status,
		       last_used_at, expires_at, created_at
		FROM api_keys`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apiKeys []*API_Keys

	for rows.Next() {
		var apiKey API_Keys

		err := rows.Scan(
			&apiKey.ID,
			&apiKey.ConsumerID,
			&apiKey.KeyHash,
			&apiKey.KeyPrefix,
			&apiKey.Status,
			&apiKey.LastUsedAt,
			&apiKey.ExpiresAt,
			&apiKey.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		apiKeys = append(apiKeys, &apiKey)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return apiKeys, nil
}
