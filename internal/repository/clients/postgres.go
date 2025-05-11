package clients

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SpaceSlow/loadbalancer/config"
	"github.com/SpaceSlow/loadbalancer/internal/domain/clients"
)

const (
	createClientQuery = `
INSERT INTO clients 
    (client_id, api_key, capacity, rps)
VALUES
    ($1, $2, $3, $4)`

	listClientsQuery = `
SELECT client_id, api_key, capacity, rps
FROM clients`

	fetchClientQuery = `
SELECT client_id, api_key, capacity, rps
FROM clients
WHERE id=$1`

	updateClientQuery = `
UPDATE clients
SET api_key=$2, capacity=$3, rps=$4
WHERE client_id=$1
`

	deleteClientQuery = `
DELETE
FROM clients
WHERE id=$1`
)

type PostgresRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresRepo(ctx context.Context, cfg config.DBConfig) (*PostgresRepo, error) {
	pool, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("NewPostgresRepo: %w", err)
	}
	return &PostgresRepo{pool: pool}, nil
}

func (r *PostgresRepo) Create(ctx context.Context, clientID, apiKey string, capacity, rps float64) (*clients.Client, error) {
	_, err := r.pool.Exec(ctx, createClientQuery, clientID, apiKey, capacity, rps)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return nil, clients.ErrClientExists
	} else if err != nil {
		return nil, err
	}

	return &clients.Client{
		ID:       clientID,
		APIKey:   apiKey,
		Capacity: capacity,
		RPS:      rps,
	}, nil
}

func (r *PostgresRepo) List(ctx context.Context) ([]clients.Client, error) {
	rows, err := r.pool.Query(ctx, listClientsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clientSlice := make([]clients.Client, 0)
	for rows.Next() {
		var client clients.Client
		if err := rows.Scan(&client.ID, &client.APIKey, &client.Capacity, &client.RPS); err != nil {
			return nil, err
		}
		clientSlice = append(clientSlice, client)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return clientSlice, nil
}

func (r *PostgresRepo) Fetch(ctx context.Context, clientID string) (*clients.Client, error) {
	row := r.pool.QueryRow(ctx, fetchClientQuery, clientID)

	var client clients.Client
	err := row.Scan(&client.ID, &client.APIKey, &client.Capacity, &client.RPS)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, clients.ErrClientNotExists
	} else if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *PostgresRepo) Update(ctx context.Context, clientID, newAPIKey string, newCapacity, newRPS float64) (*clients.Client, error) {
	commandTag, err := r.pool.Exec(ctx, updateClientQuery, clientID, newAPIKey, newCapacity, newRPS)
	if err != nil {
		return nil, err
	}
	if commandTag.RowsAffected() == 0 {
		return nil, clients.ErrClientNotExists
	}
	return &clients.Client{
		ID:       clientID,
		APIKey:   newAPIKey,
		Capacity: newCapacity,
		RPS:      newRPS,
	}, nil
}

func (r *PostgresRepo) Delete(ctx context.Context, clientID string) error {
	commandTag, err := r.pool.Exec(ctx, deleteClientQuery, clientID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return clients.ErrClientNotExists
	}
	return nil
}
