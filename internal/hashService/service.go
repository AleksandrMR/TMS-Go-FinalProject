package hashService

import (
	"encoding/hex"
	"fmt"
	"github.com/Egor-Golang-TSM-Course/final-project-AleksandrMR/internal/crypto"
	"github.com/Egor-Golang-TSM-Course/final-project-AleksandrMR/internal/domain/models"
	"golang.org/x/net/context"
	"log/slog"
)

var (
	ErrHashExists   = "hash already exists"
	ErrHashNotFound = "hash not found"
)

const (
	emptyValue = 0
)

// HashService describes the hash service
type HashService struct {
	log          *slog.Logger
	hashProvider HashProvider
}

// HashProvider contains methods that must be implemented in data storage layer
type HashProvider interface {
	CheckHashDB(ctx context.Context, hash string) (bool, error)
	GetHashDB(ctx context.Context, hash string) (models.Hash, error)
	SaveHashDB(ctx context.Context, hash string, payload string) (int64, error)
}

// New creates a new service using the passed provider
func New(
	log *slog.Logger,
	hashProvider HashProvider,
) *HashService {
	return &HashService{
		log:          log,
		hashProvider: hashProvider,
	}
}

// CheckHash queries the data storage layer for a hash
func (s *HashService) CheckHash(
	ctx context.Context,
	payload string,
) (bool, error) {
	const op = "hashService.CheckHash"
	log := s.log.With(
		slog.String("op", op),
		slog.String("payload", payload),
	)
	hash := generateHash(payload)
	log.Info("Checking Hash on DB")

	exist, err := s.hashProvider.CheckHashDB(ctx, hash)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("checked if hash is exist", slog.Bool("is_exist", exist))
	return exist, nil
}

// GetHash requests the data storage layer to obtain a hash
func (s *HashService) GetHash(
	ctx context.Context,
	payload string,
) (string, error) {
	const op = "hashService.GetHash"
	log := s.log.With(
		slog.String("op", op),
		slog.String("payload", payload),
	)
	hash := generateHash(payload)
	log.Info("Getting Hash from DB")

	dbHash, err := s.hashProvider.GetHashDB(ctx, hash)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return dbHash.HashValue, nil
}

// CreateHash requests the data storage layer to create a hash
func (s *HashService) CreateHash(
	ctx context.Context,
	payload string,
) (bool, error) {
	const op = "hashService.CreateHash"
	log := s.log.With(
		slog.String("op", op),
		slog.String("payload", payload),
	)
	hash := generateHash(payload)
	log.Info("Creating Hash on DB")

	id, err := s.hashProvider.SaveHashDB(ctx, hash, payload)
	if err != nil {
		log.Error("failed to save hash", err)
		return false, fmt.Errorf("%s: %w", op, err)
	}
	isCreated := false
	if id != emptyValue {
		isCreated = true
	}
	return isCreated, nil
}

func generateHash(payload string) string {
	hash := crypto.NewHashSHA256([]byte(payload))
	sHash := hex.EncodeToString(hash)
	return sHash
}
