package usecase

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "errors"

    "github.com/kaka/kodeakademia/be/internal/domain"
    "github.com/kaka/kodeakademia/be/internal/domain/repository"
)

// ErrNotImplemented is returned by stubs until implemented.
var ErrNotImplemented = errors.New("not implemented")

// LoginOrProvisionUser handles login or provisioning of a user from an OAuth profile.
// It returns the user, a token string, or an error.
func LoginOrProvisionUser(ctx context.Context, repo repository.UserRepository, googleProfile map[string]string) (*domain.User, string, error) {
    // Expect `sub` field in the googleProfile
    sub, ok := googleProfile["sub"]
    if !ok || sub == "" {
        return nil, "", errors.New("google profile missing sub")
    }

    // 1) Try to find existing user by Google sub
    u, err := repo.FindByGoogleSub(ctx, sub)
    if err != nil {
        return nil, "", err
    }

    if u != nil {
        // existing user: issue a simple token and return
        token, tokErr := generateDummyToken()
        if tokErr != nil {
            return nil, "", tokErr
        }
        return u, token, nil
    }

    // 2) Not found: create new user from profile
    newUser := &domain.User{
        Email:     googleProfile["email"],
        Name:      googleProfile["name"],
        AvatarURL: googleProfile["picture"],
        GoogleSub: sub,
        Role:      "learner",
    }

    created, err := repo.CreateFromOAuth(ctx, newUser)
    if err != nil {
        return nil, "", err
    }

    token, tokErr := generateDummyToken()
    if tokErr != nil {
        return nil, "", tokErr
    }

    return created, token, nil
}

func generateDummyToken() (string, error) {
    b := make([]byte, 16)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil
}