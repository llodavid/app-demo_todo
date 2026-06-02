package commons

import (
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnvFile() string {
	// external configuration using environment variables
	var name = os.Getenv("ENV_NAME")
	if len(name) == 0 {
		name = ".env"
	}
	godotenv.Load(name)
	return name
}

func parseLogLevelName(levelName string) slog.Level {
	logLevel := slog.LevelInfo
	switch strings.ToUpper(levelName) {
	case slog.LevelDebug.String():
		logLevel = slog.LevelDebug
	case slog.LevelWarn.String():
		logLevel = slog.LevelWarn
	case slog.LevelError.String():
		logLevel = slog.LevelError
	}
	return logLevel
}

func InitLoggerFromEnv() *slog.Logger {
	// init structured logging with configuration
	levelName := os.Getenv("LOG_LEVEL")
	if len(levelName) == 0 {
		levelName = "info"
	}
	sourceLocation, err := strconv.ParseBool(os.Getenv("LOG_SOURCE"))
	if err != nil {
		sourceLocation = false
	}
	logLevel := parseLogLevelName(levelName)
	opts := &slog.HandlerOptions{Level: logLevel, AddSource: sourceLocation}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}
