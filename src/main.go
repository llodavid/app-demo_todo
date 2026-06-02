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

var store *app.Store

func main() {
	commons.LoadEnvFile()
	commons.InitLoggerFromEnv()
	slog.Debug("main::main() - Executing")
	//
	store, err := app.NewStore()
	if err != nil {
		slog.Error("main::main() - Failed to create store", "error", err)
		return
	}
	defer store.Destroy()
	//
	router := http.NewServeMux()
	app.NewUi(router, store)
	//
	srv, _ := commons.NewServer(router)
	host := os.Getenv("APP_HOST")
	port := os.Getenv("APP_PORT")
	slog.Info("main::main() - Web Server is available at http://" + host + ":" + port)
	slog.Info("main::main() - Press Ctrl+C to stop")
	if err := srv.RunServer(context.Background(), 5*time.Second); err != nil {
		slog.Error("main::main() - Server error", "error", err)
	}
	slog.Info("main::main() - Stopped")
}
