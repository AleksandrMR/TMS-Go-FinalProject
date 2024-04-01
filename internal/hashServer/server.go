package hashServer

import (
	"errors"
	desc "github.com/AleksandrMR/proto_hashService/gen/hashService_v1"
	"github.com/Egor-Golang-TSM-Course/final-project-AleksandrMR/internal/storage"
	"golang.org/x/net/context"
	grps "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HashService interface for implementing business logic in the service layer
type HashService interface {
	CheckHash(ctx context.Context, payload string) (bool, error)
	GetHash(ctx context.Context, payload string) (string, error)
	CreateHash(ctx context.Context, payload string) (bool, error)
}

// serverAPI contains a server and service
type serverAPI struct {
	desc.UnsafeHashServiceServer
	hashService HashService
}

// Register registers the server and service
func Register(gRPS *grps.Server, hashService HashService) {
	desc.RegisterHashServiceServer(gRPS, &serverAPI{hashService: hashService})
}

// CheckHash endpoint "/hashService/v1/checkHash"
func (s *serverAPI) CheckHash(
	ctx context.Context,
	request *desc.CheckHashRequest,
) (*desc.CheckHashResponse, error) {
	if err := validateCheckHash(request); err != nil {
		return nil, err
	}
	exist, err := s.hashService.CheckHash(ctx, request.GetPayload())
	if err != nil {
		// TODO: handle specific error cases
		if errors.Is(err, storage.ErrHashNotFound) {
			return &desc.CheckHashResponse{
				HashExist: exist,
			}, nil
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &desc.CheckHashResponse{
		HashExist: exist,
	}, nil
}

// GetHash endpoint "/hashService/v1/getHash"
func (s *serverAPI) GetHash(
	ctx context.Context,
	request *desc.GetHashRequest,
) (*desc.GetHashResponse, error) {
	if err := validateGetHash(request); err != nil {
		return nil, err
	}
	hash, err := s.hashService.GetHash(ctx, request.GetPayload())
	if err != nil {
		// TODO: handle specific error cases
		if errors.Is(err, storage.ErrHashNotFound) {
			return nil, status.Error(codes.NotFound, "hash not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &desc.GetHashResponse{
		Hash: hash,
	}, nil
}

// CreateHash endpoint "/hashService/v1/createHash"
func (s *serverAPI) CreateHash(
	ctx context.Context,
	request *desc.CreateHashRequest,
) (*desc.CreateHashResponse, error) {
	if err := validateCreateHash(request); err != nil {
		return nil, err
	}
	isCreated, err := s.hashService.CreateHash(ctx, request.GetPayload())
	if err != nil {
		// TODO: handle specific error cases
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &desc.CreateHashResponse{
		HashCreated: isCreated,
	}, nil
}

// ---------------- Validation Functions ----------------------------------

func validateCheckHash(req *desc.CheckHashRequest) error {
	if req.GetPayload() == "" {
		return status.Error(codes.InvalidArgument, "payload is required")
	}
	return nil
}

func validateGetHash(req *desc.GetHashRequest) error {
	if req.GetPayload() == "" {
		return status.Error(codes.InvalidArgument, "payload is required")
	}
	return nil
}

func validateCreateHash(req *desc.CreateHashRequest) error {
	if req.GetPayload() == "" {
		return status.Error(codes.InvalidArgument, "payload is required")
	}
	return nil
}
