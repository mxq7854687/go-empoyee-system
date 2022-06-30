package db

import (
	"context"
	"example/employee/server/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
		Status:         UserStatusPending,
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.Email, arg.Email)
	require.Equal(t, user.HashedPassword, arg.HashedPassword)
	require.Equal(t, user.Status, arg.Status)

	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)

	require.WithinDuration(t, user1.UpdatedAt, user2.UpdatedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestResetUserPassword(t *testing.T) {
	user1 := createRandomUser(t)
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	// TODO add hook to update the time when trigger the query instead of manually update
	arg := ActivateUserParams{
		Email:          user1.Email,
		HashedPassword: hashedPassword,
	}

	testQueries.ActivateUser(context.TODO(), arg)

	updatedUser, err := testQueries.GetUser(context.Background(), user1.Email)
	require.NoError(t, err)
	require.Equal(t, hashedPassword, updatedUser.HashedPassword)
	require.True(t, updatedUser.UpdatedAt.After(user1.UpdatedAt))
	require.Equal(t, updatedUser.Status, UserStatusActivated)
}
