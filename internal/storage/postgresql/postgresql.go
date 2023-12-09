package postgresql

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, storagePath string) (*Storage, error) {
	const op = "storage.postgresql.New"

	dbpool, err := pgxpool.New(ctx, storagePath)
	if err != nil {
		os.Exit(1)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	prepareQuery := `
		CREATE TABLE IF NOT EXISTS url(
        id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
        alias TEXT NOT NULL UNIQUE,
        url TEXT NOT NULL);
    CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`

	_, err = dbpool.Exec(ctx, prepareQuery)

	if err != nil {
		os.Exit(1)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: dbpool,
	}, nil
}

func (storage *Storage) Ping(ctx context.Context) error {
	return storage.db.Ping(ctx)
}

func (storage *Storage) Close() {
	storage.db.Close()
}

func (storage *Storage) SaveURL(ctx context.Context, urlToSave string, alias string) (string, error) {
	const op = "storage.postgresql.SaveURL"

	query := `INSERT INTO url (url,alias) VALUES (@urlToSave, @alias) RETURNING id;`
	args := pgx.NamedArgs{
		"urlToSave": urlToSave,
		"alias":     alias,
	}

	var id string
	err := storage.db.QueryRow(ctx, query, args).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("unable to insert row: %w", err)
	}

	return id, nil
}

func (storage *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "storage.postgresql.GetURL"

	query := `SELECT url FROM url WHERE alias=(@alias);`
	args := pgx.NamedArgs{
		"alias": alias,
	}

	var url string
	err := storage.db.QueryRow(ctx, query, args).Scan(&url)

	if err != nil {
		return "", fmt.Errorf("unable to get row: %w", err)
	}

	return url, nil
}

func (storage *Storage) DeleteURL(ctx context.Context, id string) error {
	const op = "storage.postgresql.DeleteURL"

	query := `DELETE FROM url WHERE id=(@id);`
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := storage.db.Exec(ctx, query, args)

	if err != nil {
		return fmt.Errorf("unable to delete row: %w", err)
	}

	return nil
}
