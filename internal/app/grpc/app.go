package grpcapp

import (
	"fmt"
	"github.com/Egor-Golang-TSM-Course/final-project-AleksandrMR/internal/hashServer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPSServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	hashService hashServer.HashService,
	port int,
) *App {
	grpcServer := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
	)
	reflection.Register(grpcServer)
	hashServer.Register(grpcServer, hashService)

	return &App{
		log:        log,
		gRPSServer: grpcServer,
		port:       port,
	}
}

func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (app *App) Run() error {
	const op = "grpcapp.Run"
	log := app.log.With(
		slog.String("op", op),
		slog.Int("port", app.port),
	)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("gRPC server is running", slog.String("addr", l.Addr().String()))

	if err := app.gRPSServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (app *App) Stop() {
	const op = "grpcapp.Stop"
	app.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", app.port))
	app.gRPSServer.GracefulStop()
}
