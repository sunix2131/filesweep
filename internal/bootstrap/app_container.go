package bootstrap

import (
	"database/sql"
	"io"
	"log/slog"

	"filesweep/internal/application"
	"filesweep/internal/infrastructure/apppaths"
	"filesweep/internal/infrastructure/database/repositories"
)

type Container struct {
	Paths    apppaths.Paths
	DB       *sql.DB
	Logger   *slog.Logger
	logFile  io.Closer
	Services *application.Services
}

func NewContainer() (*Container, error) {
	paths, err := apppaths.Resolve()
	if err != nil {
		return nil, err
	}
	logger, logFile, err := NewLogger(paths)
	if err != nil {
		return nil, err
	}
	db, err := OpenDatabase(paths.DBPath)
	if err != nil {
		_ = logFile.Close()
		return nil, err
	}
	store := repositories.NewStore(db)
	c := &Container{Paths: paths, DB: db, Logger: logger, logFile: logFile}
	c.Services = application.NewServices(store, paths, logger, nil)
	return c, nil
}

func (c *Container) Close() error {
	if c.DB != nil {
		_ = c.DB.Close()
	}
	if c.logFile != nil {
		return c.logFile.Close()
	}
	return nil
}
