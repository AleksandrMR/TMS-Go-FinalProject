package app

import (
	grpcapp "github.com/Egor-Golang-TSM-Course/final-project-AleksandrMR/internal/app/grpc"
	"github.com/Egor-Golang-TSM-Course/final-project-AleksandrMR/internal/hashService"
	"github.com/Egor-Golang-TSM-Course/final-project-AleksandrMR/internal/storage/sqlite"
	"log/slog"
)

type App struct {
	GRPCServ *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}
	service := hashService.New(log, storage)
	grpcNewApp := grpcapp.New(log, service, grpcPort)

	return &App{
		GRPCServ: grpcNewApp,
	}
}
