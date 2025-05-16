package integration

import (
	"context"
	"testing"
	"time"

	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthRepository_Create(t *testing.T) {
	// Arrange
	ctx := context.Background()
	phoneNumber := "09123456789"
	user := &entities.User{
		ID:          uuid.New(),
		PhoneNumber: phoneNumber,
		Password:    "hashed_password",
		Status:      entities.Active,
		Role:        entities.UserRole,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Act
	err := authRepo.Create(ctx, user)

	// Assert
	require.NoError(t, err)

	// Cleanup
	err = cleanupDatabase()
	require.NoError(t, err)
}

func TestAuthRepository_FindUserByPhoneNumber(t *testing.T) {
	// Arrange
	ctx := context.Background()
	phoneNumber := "09123456789"
	user := &entities.User{
		ID:          uuid.New(),
		PhoneNumber: phoneNumber,
		Password:    "hashed_password",
		FirstName:   "Test",
		LastName:    "User",
		Email:       "test@example.com",
		Status:      entities.Active,
		Role:        entities.UserRole,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := authRepo.Create(ctx, user)
	require.NoError(t, err)

	// Act
	foundUser, err := authRepo.FindUserByPhoneNumber(ctx, &phoneNumber)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, phoneNumber, foundUser.PhoneNumber)
	assert.Equal(t, user.ID, foundUser.ID)

	// Cleanup
	err = cleanupDatabase()
	require.NoError(t, err)
}

func TestAuthRepository_FindUserByID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	user := &entities.User{
		ID:          uuid.New(),
		PhoneNumber: "09123456789",
		Password:    "hashed_password",
		FirstName:   "Test",
		LastName:    "User",
		Email:       "test@example.com",
		Status:      entities.Active,
		Role:        entities.UserRole,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := authRepo.Create(ctx, user)
	require.NoError(t, err)

	// Act
	foundUser, err := authRepo.FindUserByID(ctx, user.ID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.PhoneNumber, foundUser.PhoneNumber)

	// Cleanup
	err = cleanupDatabase()
	require.NoError(t, err)
}

func TestAuthRepository_FindUserByPhoneNumber_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	phoneNumber := "09999999999"

	// Act
	user, err := authRepo.FindUserByPhoneNumber(ctx, &phoneNumber)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, errors.ErrUserNotFound, err)

	// Cleanup
	err = cleanupDatabase()
	require.NoError(t, err)
}

func TestAuthRepository_FindUserByID_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	id := uuid.New()

	// Act
	user, err := authRepo.FindUserByID(ctx, id)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, errors.ErrUserNotFound, err)

	// Cleanup
	err = cleanupDatabase()
	require.NoError(t, err)
}
