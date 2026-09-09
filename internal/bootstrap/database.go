package bootstrap

import (
	"database/sql"
	"net/url"
	"path/filepath"
	"strings"

	dbinfra "filesweep/internal/infrastructure/database"

	_ "modernc.org/sqlite"
)

func OpenDatabase(path string) (*sql.DB, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	uriPath := filepath.ToSlash(absolute)
	if !strings.HasPrefix(uriPath, "/") {
		uriPath = "/" + uriPath
	}
	// PRAGMAs are applied by the driver to every connection, including replacements.
	dsn := url.URL{Scheme: "file", Path: uriPath, RawQuery: "_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)"}
	db, err := sql.Open("sqlite", dsn.String())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(`PRAGMA foreign_keys = ON; PRAGMA journal_mode = WAL; PRAGMA busy_timeout = 5000;`); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := dbinfra.RunMigrations(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
