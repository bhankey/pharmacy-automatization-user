package app

import (
	"context"
	"fmt"
	"github.com/bhankey/pharmacy-automatization-user/pkg/api/userservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bhankey/pharmacy-automatization-user/internal/app/container"
	configinternal "github.com/bhankey/pharmacy-automatization-user/internal/config"
	"github.com/bhankey/pharmacy-automatization-user/pkg/logger"
)

type App struct {
	server    *grpc.Server
	listener  net.Listener
	container *container.Container
	logger    logger.Logger
}

const shutDownTimeoutSeconds = 10

func NewApp(configPath string) (*App, error) {
	config, err := configinternal.GetConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to init app because of config error: %w", err)
	}

	log, err := logger.GetLogger(config.Logger.Path, config.Logger.Level, true)
	if err != nil {
		return nil, fmt.Errorf("failed to init logger error: %w", err)
	}

	log.Info("try to init data source resource")
	dataSources, err := newDataSource(config) // TODO remove dataSource struct
	if err != nil {
		return nil, err
	}

	smtp, err := newSMTPClient(config)
	if err != nil {
		return nil, err
	}

	dependencies := container.NewContainer(
		log,
		dataSources.db,
		dataSources.db,
		dataSources.redisClient,
		smtp,
		config.Secure.JwtKey,
		config.SMTP.From,
	)

	grpcHandler := dependencies.GetUserGRPCHandler()

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	userservice.RegisterUserServiceServer(grpcServer, grpcHandler)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.Server.Port))
	if err != nil {
		return nil, err
	}

	return &App{
		server:    grpcServer,
		listener:  listener,
		container: dependencies,
		logger:    log,
	}, nil
}

func (a *App) Start() {
	a.logger.Info("staring server on addr: " + a.listener.Addr().String())
	go func() {

		if err := a.server.Serve(a.listener); err != nil {
			a.logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	a.logger.Info("received signal to shutdown server")
	ctx, cancel := context.WithTimeout(context.Background(), shutDownTimeoutSeconds*time.Second)
	defer cancel()
	a.server.GracefulStop()

	<-ctx.Done()

	a.container.CloseAllConnections()

	a.logger.Info("server was shutdown")
}
