package application

import (
	"context"
	"fmt"
	"log"
	"newsfeed/internal/adapters/http"
	"newsfeed/internal/adapters/repository"
	"newsfeed/internal/config"
	"newsfeed/pkg/httpServer"
	"newsfeed/pkg/logging"
	"newsfeed/pkg/mongoDB"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *config.Config) {
	logger, err := logging.New(cfg.Logger.Level, cfg.Logger.Format)
	if err != nil {
		log.Fatal(err)
	}
	logger.Debug("logger created")

	logger.Debug("connecting to the database")

	mongodb, err := mongoDB.New(logger, cfg.Database.Host, cfg.Database.Port, 3)
	if err != nil {
		logger.Fatalf("failed to connect to the databaseЖ %s", err.Error())
	}
	repo, err := repository.New(cfg, mongodb)

	logger.Debug("successfully")

	logger.Debug("creating instance of rest service")
	httpHandler := http.NewHandler(repo)
	restServer := httpServer.New(
		httpHandler.GetHandler(),
		cfg.Http.Port,
		cfg.Http.ReadTimeout,
		cfg.Http.WriteTimeout,
	)
	logger.Debug("successfully, running it")
	restServer.Run()

	// Shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-interrupt:
		logger.Info("an interrupt signal was received " + s.String())
	case err = <-restServer.Notify():
		logger.Error(fmt.Errorf("httpServer.Notify: %w", err))
	}

	ctx, cancelFn := context.WithTimeout(
		context.Background(),
		time.Duration(cfg.HyperParams.ShutdownTimeout)*time.Second,
	)
	defer cancelFn()
	if err := restServer.Shutdown(ctx); err != nil {
		logger.Errorf("error during shutdown httpServer: %s", err.Error())
	}
	if err := mongodb.Shutdown(ctx); err != nil {
		logger.Errorf("error during shutdown DB: %s", err.Error())
	}

}
