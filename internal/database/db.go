package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/manjunath-tintbytes/seafoodai.api/internal/migrations"
	"github.com/manjunath-tintbytes/seafoodai.api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var sqlDB *sql.DB

func SetupDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s options='-c client_encoding=UTF8'",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("SSL_MODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to DB: ", err)
	}

	sqlDB, err = db.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB: ", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Create users table
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	sqlDB.Exec(createUsersTable)

	// Create password_reset_tokens table
	createTokensTable := `
	CREATE TABLE IF NOT EXISTS password_reset_tokens (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token VARCHAR(255) UNIQUE NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		used BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	sqlDB.Exec(createTokensTable)

	return db
}

// RunMigrations runs all pending migrations
func RunMigrations(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations.GetMigrations())

	if err := m.Migrate(); err != nil {
		log.Printf("Could not migrate: %v", err)
		return err
	}

	log.Printf("Migration did run successfully")
	return nil
}

// RollbackMigration rolls back the last migration
func RollbackMigration(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations.GetMigrations())

	if err := m.RollbackLast(); err != nil {
		log.Printf("Could not rollback: %v", err)
		return err
	}

	log.Printf("Rollback did run successfully")
	return nil
}

// RollbackToMigration rolls back to a specific migration
func RollbackToMigration(db *gorm.DB, migrationID string) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations.GetMigrations())

	if err := m.RollbackTo(migrationID); err != nil {
		log.Printf("Could not rollback to migration %s: %v", migrationID, err)
		return err
	}

	log.Printf("Rollback to migration %s did run successfully", migrationID)
	return nil
}

func CloseDB() {
	if sqlDB != nil {
		sqlDB.Close()
	}
}

func CreateUser(user *models.User) error {
	query := `INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id, created_at`
	return sqlDB.QueryRow(query, user.Email, user.Password, user.Name).Scan(&user.ID, &user.CreatedAt)
}

func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password, name, created_at FROM users WHERE email = $1`
	err := sqlDB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password, name, created_at FROM users WHERE id = $1`
	err := sqlDB.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CreatePasswordResetToken(userID int, token string, expiresAt time.Time) error {
	query := `INSERT INTO password_reset_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := sqlDB.Exec(query, userID, token, expiresAt)
	return err
}

func GetPasswordResetToken(token string) (*models.PasswordResetToken, error) {
	resetToken := &models.PasswordResetToken{}
	query := `SELECT id, user_id, token, expires_at, used, created_at FROM password_reset_tokens WHERE token = $1`
	err := sqlDB.QueryRow(query, token).Scan(
		&resetToken.ID,
		&resetToken.UserID,
		&resetToken.Token,
		&resetToken.ExpiresAt,
		&resetToken.Used,
		&resetToken.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return resetToken, nil
}

func MarkTokenAsUsed(token string) error {
	query := `UPDATE password_reset_tokens SET used = true WHERE token = $1`
	_, err := sqlDB.Exec(query, token)
	return err
}

func UpdateUserPassword(userID int, hashedPassword string) error {
	query := `UPDATE users SET password = $1 WHERE id = $2`
	_, err := sqlDB.Exec(query, hashedPassword, userID)
	return err
}
