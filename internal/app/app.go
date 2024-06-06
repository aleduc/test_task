// Package app implements methods for application start/shutdown.
package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gorm.io/driver/postgres"

	"test_task/internal/app/server"
	"test_task/internal/config"
	httpController "test_task/internal/controller/http"
	"test_task/internal/datastore/kafka"
	postgresRepo "test_task/internal/datastore/postgres"
	"test_task/internal/notificator"
	"test_task/internal/usecase"
)

// Start service.
func Start(cfg config.Config) {
	mainCtx, mainCtxCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer mainCtxCancel()

	l := logrus.New()
	pgClient, err := gorm.Open(postgres.Open(cfg.Postgres.URL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error)})
	if err != nil {
		l.Fatalf("can't connect to the Postgres database: %s", err.Error())
		return
	}

	notificatorCtx, notificatorCtxCancel := context.WithCancel(context.Background())
	defer notificatorCtxCancel()
	notitifcationPubSub := notificator.NewPubSub(&kafka.Mock{}, cfg.Notification.BufferSize, l)
	go func() {
		notitifcationPubSub.Start(notificatorCtx)
	}()

	userRepo := postgresRepo.NewUserRepository(pgClient)
	userUseCase := usecase.NewUser(userRepo, notitifcationPubSub, l)
	userController := httpController.NewUserHandler(userUseCase, l)

	healthRepo := postgresRepo.NewHealth(pgClient)
	healthController := httpController.NewHealthController([]httpController.HealthChecker{healthRepo})
	echoServer := server.NewServer(cfg.HTTP)
	httpController.InitRoutes(echoServer, httpController.Controllers{User: userController, HealthController: healthController})

	serverStopped := make(chan struct{}, 1)
	go func() {
		l.Info("HTTP server is started")
		if err := echoServer.Start(fmt.Sprintf(":%s", cfg.HTTP.Port)); !errors.Is(err, http.ErrServerClosed) {
			l.Errorf("shutting down http server: %s", err.Error())
		}
		serverStopped <- struct{}{}
	}()

	select {
	case <-mainCtx.Done():
		l.Info("graceful shutting down…")
		server.ShutdownServer(l, echoServer, cfg.HTTP)
		<-serverStopped
		l.Info("http server is stopped")
		notitifcationPubSub.Stop(notificatorCtxCancel, cfg.Notification.RecheckTimeout, cfg.Notification.CloseTimeout)
	case <-serverStopped:
		server.ShutdownServer(l, echoServer, cfg.HTTP)
		l.Info("http server is stopped")
		notitifcationPubSub.Stop(notificatorCtxCancel, cfg.Notification.RecheckTimeout, cfg.Notification.CloseTimeout)
	}

	l.Info("service is stopped")
}
