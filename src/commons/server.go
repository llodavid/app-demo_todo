package commons

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	hs *http.Server
}

func NewServer(router *http.ServeMux, port string) *Server {
	fs := http.FileServer(http.Dir("./public"))
	router.Handle("GET /public/", http.StripPrefix("/public/", fs))
	//
	intPort := port
	if IsRunningInDockerContainer() {
		// internal port is always 80; see dockerfile for port mapping
		intPort = "80"
	}
	srv := http.Server{
		Addr:    ":" + intPort,
		Handler: router,
	}
	return &Server{hs: &srv}
}

func (s *Server) RunServer(
	ctx context.Context,
	shutdownTimeout time.Duration,
) error {
	// see [Graceful Shutdown in Go: Key Patterns you need to know!](https://www.youtube.com/watch?v=UPVSeZXBTxI)
	serverErr := make(chan error, 1)
	go func() {
		slog.Debug("server::RunServer() - Started listening")
		if err := s.hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
		close(serverErr)
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	select {
	case err := <-serverErr:
		return err
	case <-stop:
		slog.Debug("server::RunServer() - Shutdown signal received, stopping server")
	case <-ctx.Done():
		slog.Debug("server::RunServer() - Context cancelled, stopping server")
	}
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		shutdownTimeout)
	defer cancel()
	if err := s.hs.Shutdown(shutdownCtx); err != nil {
		if closeErr := s.hs.Close(); closeErr != nil {
			return errors.Join(err, closeErr)
		}
		return err
	}
	slog.Debug("server::RunServer() - Server stopped gracefully")
	return nil
}

func IsRunningInDockerContainer() bool {
	// docker creates a .dockerenv file at the root
	// of the directory tree inside the container.
	// if this file exists then the viewer is running
	// from inside a container so return true
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}
