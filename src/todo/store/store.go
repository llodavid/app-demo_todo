package todo_store

import (
	"RobertTC32/example-demo_hello/src/commons"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// "store" (aka repository) contains adapter for MariaDB DB in the "Hexagonal Architecture"

type Store struct {
	Dsn string
	Db  *sqlx.DB
}

func NewStore() (*Store, error) {
	slog.Debug("store::NewStore() - Executing")
	//dsnFormat := "%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local"
	//dsn := fmt.Sprintf(dsnFormat, Username, Password, Host, Port, DbName)
	dbdriver := os.Getenv("DB_DRIVER")
	dsn := os.Getenv("DB_DSN")
	if commons.IsRunningInDockerContainer() {
		// change database localhost to host.docker.internal
		dsn = strings.Replace(dsn, "localhost", "host.docker.internal", 1)
	}
	db, err := sqlx.Open(dbdriver, dsn)
	if err != nil {
		slog.Error("store::NewStore() - Failed to open SQL connection", "error", err)
		return &Store{}, fmt.Errorf("Failed to open SQL connection: %w", err)
	}
	// Pool configuration
	poolFlag, err := strconv.Atoi(os.Getenv("DB_POOL_MAX_OPEN_CONNS"))
	if err == nil {
		slog.Debug("store::NewStore() - Set", "maxOpenConns", poolFlag)
		db.SetMaxOpenConns(poolFlag) // max simultaneous connections
	}
	poolFlag, err = strconv.Atoi(os.Getenv("DB_POOL_MAX_IDLE_CONNS"))
	if err == nil {
		slog.Debug("store::NewStore() - Set", "maxIdleConns", poolFlag)
		db.SetMaxIdleConns(poolFlag) // connections kept idle in pool
	}
	poolFlag, err = strconv.Atoi(os.Getenv("DB_POOL_CONN_MAX_LIFETIME"))
	if err == nil {
		slog.Debug("store::NewStore() - Set", "connMaxLifetime", poolFlag)
		db.SetConnMaxLifetime(time.Duration(poolFlag) * time.Minute) // max connection age
	}
	poolFlag, err = strconv.Atoi(os.Getenv("DB_POOL_CONN_MAX_IDLE_TIME"))
	if err == nil {
		slog.Debug("store::NewStore() - Set", "connMaxIdleTime", poolFlag)
		db.SetConnMaxIdleTime(time.Duration(poolFlag) * time.Minute) // evict idle connections after this
	}
	err = db.Ping()
	if err != nil {
		slog.Error("store::NewStore() - Failed to ping db", "error", err)
		return &Store{}, fmt.Errorf("Failed to ping db: %w", err)
	}
	slog.Debug("store::NewStore() - Database opened", "dsn", dsn)
	return &Store{
		Dsn: dsn,
		Db:  db,
	}, nil
}

func (s *Store) Destroy() error {
	slog.Debug("store::Destroy() - Executing")
	return s.Db.Close()
}
