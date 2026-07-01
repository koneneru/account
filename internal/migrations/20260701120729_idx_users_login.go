package migrations

import (
	"context"
	"database/sql"
)

func upIdxUsersLogin(ctx context.Context, tx *sql.Tx) error {
	query := `CREATE INDEX IF NOT EXISTS idx_users_login ON users(login);`
	_, err := tx.Exec(query)
	return err
}

func downIdxUsersLogin(ctx context.Context, tx *sql.Tx) error {
	query := "DROP INDEX IF EXISTS idx_users_login ON users;"
	_, err := tx.Exec(query)
	return err
}
