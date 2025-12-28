package usecase

import (
    "context"
    "testing"

    "github.com/stretchr/testify/require"

    "github.com/kaka/kodeakademia/be/internal/domain"
)

// mockRepo is a simple in-memory mock implementing UserRepository for tests.
type mockRepo struct{
    found *domain.User
}

func (m *mockRepo) FindByGoogleSub(ctx context.Context, googleSub string) (*domain.User, error) {
    if m.found != nil && m.found.GoogleSub == googleSub {
        return m.found, nil
    }
    return nil, nil
}

func (m *mockRepo) CreateFromOAuth(ctx context.Context, u *domain.User) (*domain.User, error) {
    // simulate assigning ID
    u.ID = 123
    return u, nil
}

func (m *mockRepo) Update(ctx context.Context, u *domain.User) (*domain.User, error) {
    return u, nil
}

func TestLoginOrProvisionUser_ExistingUser(t *testing.T) {
    // Red: write failing test first that expects an existing user to be returned with a token
    ctx := context.Background()
    existing := &domain.User{ID: 1, Email: "student1@example.com", Name: "Student One", GoogleSub: "student-sub-1", Role: "learner"}
    mr := &mockRepo{found: existing}

    profile := map[string]string{"sub": "student-sub-1", "email": "student1@example.com", "name": "Student One"}

    user, token, err := LoginOrProvisionUser(ctx, mr, profile)

    // We expect no error and a non-empty token (this will fail until implementation is provided)
    require.NoError(t, err)
    require.NotNil(t, user)
    require.NotEmpty(t, token)
}
