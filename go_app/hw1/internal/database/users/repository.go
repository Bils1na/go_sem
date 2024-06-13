package users

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"gitlab.com/robotomize/gb-golang/homework/03-01-umanager/internal/database"
)

func New(userDB *pgx.Conn, timeout time.Duration) *Repository {
	return &Repository{userDB: userDB, timeout: timeout}
}

type Repository struct {
	userDB  *pgx.Conn
	timeout time.Duration
}

func (r *Repository) Create(ctx context.Context, req CreateUserReq) (database.User, error) {
	var u database.User

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `INSERT INTO users (id, username, password, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING id, username, created_at, updated_at`

	err := r.userDB.QueryRow(ctx, query, req.ID, req.Username, req.Password, time.Now(), time.Now()).Scan(&u.ID, &u.Username, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return database.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return u, nil
}

func (r *Repository) FindByID(ctx context.Context, userID uuid.UUID) (database.User, error) {
	var u database.User

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `SELECT id, username, created_at, updated_at
			  FROM users
			  WHERE id = $1`

	err := r.userDB.QueryRow(ctx, query, userID).Scan(&u.ID, &u.Username, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return database.User{}, fmt.Errorf("user not found")
		}
		return database.User{}, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return u, nil
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (database.User, error) {
	var u database.User

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `SELECT id, username, created_at, updated_at
			  FROM users
			  WHERE username = $1`

	err := r.userDB.QueryRow(ctx, query, username).Scan(&u.ID, &u.Username, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return database.User{}, fmt.Errorf("user not found")
		}
		return database.User{}, fmt.Errorf("failed to find user by username: %w", err)
	}

	return u, nil
}
