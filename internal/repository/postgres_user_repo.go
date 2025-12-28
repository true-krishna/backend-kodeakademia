package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kaka/kodeakademia/be/internal/domain"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) FindByGoogleSub(ctx context.Context, googleSub string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, email, name, avatar_url, google_sub, role FROM users WHERE google_sub = $1`, googleSub)
	var u domain.User
	if err := row.Scan(&u.ID, &u.Email, &u.Name, &u.AvatarURL, &u.GoogleSub, &u.Role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepository) CreateFromOAuth(ctx context.Context, u *domain.User) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO users (email, name, avatar_url, google_sub, role)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (google_sub) DO UPDATE SET email=EXCLUDED.email, name=EXCLUDED.name, avatar_url=EXCLUDED.avatar_url, role=EXCLUDED.role
		RETURNING id, email, name, avatar_url, google_sub, role
	`, u.Email, u.Name, u.AvatarURL, u.GoogleSub, u.Role)
	var out domain.User
	if err := row.Scan(&out.ID, &out.Email, &out.Name, &out.AvatarURL, &out.GoogleSub, &out.Role); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, u *domain.User) (*domain.User, error) {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET email=$1, name=$2, avatar_url=$3, role=$4 WHERE id=$5`, u.Email, u.Name, u.AvatarURL, u.Role, u.ID)
	if err != nil {
		return nil, err
	}
	return u, nil
}
