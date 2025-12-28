package repository

import (
    "context"

    "github.com/kaka/kodeakademia/be/internal/domain"
)

type UserRepository interface {
    // FindByGoogleSub returns a user by Google subject (sub) or nil if not found.
    FindByGoogleSub(ctx context.Context, googleSub string) (*domain.User, error)

    // CreateFromOAuth creates a new user record from an OAuth profile.
    CreateFromOAuth(ctx context.Context, u *domain.User) (*domain.User, error)

    // Update updates an existing user record.
    Update(ctx context.Context, u *domain.User) (*domain.User, error)
}
