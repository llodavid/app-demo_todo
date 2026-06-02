package main

import (
	"RobertTC32/example-demo_hello/src/app"
	"RobertTC32/example-demo_hello/src/commons"
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var storage *app.Storage

func main() {
	commons.LoadEnvFile()
	commons.InitLoggerFromEnv()
	slog.Debug("main::main() - Started")
	//
	dsn := os.Getenv("DB_DSN")
	storageVal, err := app.NewStorage(dsn)
	if err != nil {
		slog.Error("main::main() - Failed to create storage", "error", err)
		return
	}
	storage = &storageVal
	defer storage.Destroy()
	//
	router := http.NewServeMux()
	_, err = app.NewUi(router, storage)
	if err != nil {
		slog.Error("main::main() - Failed to create ui", "error", err)
		return
	}
	//
	port := os.Getenv("APP_PORT")
	srv := commons.NewServer(router, port)
	slog.Info("main::main() - Web Server is available at http://localhost:" + port)
	slog.Info("main::main() - Press Ctrl+C to stop")
	if err := srv.RunServer(context.Background(), 5*time.Second); err != nil {
		slog.Error("main::main() - Server error", "error", err)
	}
	slog.Info("main::main() - Stopped")
}
