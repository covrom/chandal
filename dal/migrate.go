package dal

import "github.com/jmoiron/sqlx"

func MigrateDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users
		(
			id uuid PRIMARY KEY,
			name varchar
		)`); err != nil {
		return nil, err
	}
	return db, nil
}
