// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/kordape/gymshark-task/internal/config"
	"github.com/kordape/gymshark-task/internal/httpd"
	"github.com/kordape/gymshark-task/internal/packs"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyTime:  "log_time",
		},
	})
	logrus.SetLevel(logrus.InfoLevel)

	cfg := config.Load()

	pm := packs.NewManager()

	s, err := httpd.NewServer(
		pm,
		httpd.WithPort(cfg.Port),
	)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create server")
	}

	go func() {
		if err := s.Start(); err != nil {
			logrus.WithError(err).Error("Failed to start server")
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	logrus.Info("Stopping server")
	if err := s.GracefulStop(context.Background()); err != nil {
		logrus.WithError(err).Fatal("Failed to gracefully stop the http server")
	}
}
