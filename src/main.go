package main

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var storage *Storage

func main() {
	LoadEnvFile()
	InitLoggerFromEnv()
	slog.Debug("main::main() - Started")
	//
	config := StorageConfig{
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		Host:     os.Getenv("DB_HOST"),
	}
	slog.Debug("main::main() - Storage config", "config", config)
	//
	storageVal, err := NewStorage(config)
	if err != nil {
		slog.Error("main::main() - Failed to create storage", "error", err)
		return
	}
	storage = &storageVal
	defer storage.Destroy()
	//
	port := os.Getenv("OCI_PORT")
	intPort := os.Getenv("OCI_INT_PORT")
	if !IsRunningInDockerContainer() {
		intPort = port
	}
	srv := NewServer(port, intPort)
	if err := srv.RunServer(context.Background(), 5*time.Second); err != nil {
		slog.Error("main::main() - Server error", "error", err)
	}
}

func listHandleFunc(w http.ResponseWriter, r *http.Request) {
	slog.Debug("main::listHandleFunc() - Started")
	slog.Debug("main::listHandleFunc() - Reading data", "storage", storage.config.DbName)
	d, err := storage.FindAllTodos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//
	slog.Debug("main::listHandleFunc() - Building response")
	tmpl := template.Must(template.ParseFiles("./resources/list.html"))
	err = tmpl.Execute(w, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
