package database

import "database/sql"

func Migrate(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255),
		is_verified BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS verification_codes (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) NOT NULL,
		code VARCHAR(6) NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		used BOOLEAN DEFAULT FALSE,
		verified_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Add verified_at column if it doesn't exist (for existing databases)
	ALTER TABLE verification_codes ADD COLUMN IF NOT EXISTS verified_at TIMESTAMP;

	CREATE INDEX IF NOT EXISTS idx_verification_codes_email ON verification_codes(email);
	CREATE INDEX IF NOT EXISTS idx_verification_codes_code ON verification_codes(code);
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`

	_, err := db.Exec(query)
	return err
}
