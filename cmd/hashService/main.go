package main

import (
	"github.com/Egor-Golang-TSM-Course/final-project-AleksandrMR/internal/config"
	"github.com/Egor-Golang-TSM-Course/final-project-AleksandrMR/internal/hashServer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

const (
	grpcAddress = "localhost:50051"
)

func main() {
	conf := config.MustLoad()
	log := setupLogger(conf.Env)
	log.Info("starting gRPC Server", slog.Any("conf", conf))

	if err := startGrpcServer(); err != nil {
		log.Info("startGrpcServer Error", err)
	}
}

func startGrpcServer() error {
	grpcServer := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
	)

	reflection.Register(grpcServer)
	hashServer.Register(grpcServer)

	list, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		return err
	}

	//log.Info("gRPC server listening at %v\n", grpcAddress)

	return grpcServer.Serve(list)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
