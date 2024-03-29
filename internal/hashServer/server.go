package hashServer

import (
	desc "github.com/AleksandrMR/proto_hashService/gen/hashService_v1"
	"golang.org/x/net/context"
	grps "google.golang.org/grpc"
)

type serverAPI struct {
	desc.UnsafeHashServiceServer
}

func Register(gRPS *grps.Server) {
	desc.RegisterHashServiceServer(gRPS, &serverAPI{})
}

// CheckHash endpoint "/hashService/v1/checkHash"
func (s *serverAPI) CheckHash(
	ctx context.Context,
	request *desc.CheckHashRequest,
) (*desc.CheckHashResponse, error) {
	return &desc.CheckHashResponse{
		HashExist: false,
	}, nil
}

// GetHash endpoint "/hashService/v1/getHash"
func (s *serverAPI) GetHash(
	ctx context.Context,
	request *desc.GetHashRequest,
) (*desc.GetHashResponse, error) {
	return &desc.GetHashResponse{
		Hash: "1234",
	}, nil
}

// CreateHash endpoint "/hashService/v1/createHash"
func (s *serverAPI) CreateHash(
	ctx context.Context,
	request *desc.CreateHashRequest,
) (*desc.CreateHashResponse, error) {
	return &desc.CreateHashResponse{
		HashCreated: true,
	}, nil
}
