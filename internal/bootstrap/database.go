package bootstrap

import (
	"database/sql"

	dbinfra "filesweep/internal/infrastructure/database"

	_ "modernc.org/sqlite"
)

func OpenDatabase(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
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
