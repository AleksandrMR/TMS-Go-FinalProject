package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Egor-Golang-TSM-Course/final-project-AleksandrMR/internal/domain/models"
	"github.com/Egor-Golang-TSM-Course/final-project-AleksandrMR/internal/storage"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/net/context"
)

// Storage describes the database
type Storage struct {
	db *sql.DB
}

// New create a database along the given storage path
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

// Stop closes the database
func (s *Storage) Stop() error {
	return s.db.Close()
}

// CheckHashDB check hash on DB
func (s *Storage) CheckHashDB(
	ctx context.Context,
	hash string,
) (bool, error) {
	const op = "storage.sqlite.CheckHashDB"
	_, err := s.getHashFromDB(ctx, hash, op)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetHashDB get hash from DB
func (s *Storage) GetHashDB(
	ctx context.Context,
	hash string,
) (models.Hash, error) {
	const op = "storage.sqlite.GetHashDB"
	hashModel, err := s.getHashFromDB(ctx, hash, op)
	if err != nil {
		return models.Hash{}, err
	}
	return hashModel, nil
}

// SaveHashDB save hash on DB
func (s *Storage) SaveHashDB(
	ctx context.Context,
	hash string,
	payload string,
) (int64, error) {
	const op = "storage.sqlite.SaveHashDB"
	stmt, err := s.db.Prepare("INSERT INTO hashService(hash, payload) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.ExecContext(ctx, hash, payload)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrHashExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) getHashFromDB(
	ctx context.Context,
	hash string,
	op string,
) (models.Hash, error) {
	stmt, err := s.db.Prepare("SELECT payload, hash FROM hashService WHERE hash = ?")
	if err != nil {
		return models.Hash{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, hash)
	var hashModel models.Hash
	err = row.Scan(&hashModel.Payload, &hashModel.HashValue)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Hash{}, fmt.Errorf("%s: %w", op, storage.ErrHashNotFound)
		}
		return models.Hash{}, fmt.Errorf("%s: %w", op, err)
	}
	return hashModel, nil
}
