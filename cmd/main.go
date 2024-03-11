
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mayu13/gymshark-assignment/internal/calculate"
	"github.com/mayu13/gymshark-assignment/internal/config"
	"github.com/mayu13/gymshark-assignment/internal/packs"
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

	s, err := calculate.NewServer(
		pm,
		calculate.WithPort(cfg.Port),
	)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create server")
	}

	// Create a file server to serve static files from the "web" directory
	http.Handle("/", http.FileServer(http.Dir("web")))

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
