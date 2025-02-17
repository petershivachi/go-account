package db

import (
	"context"
	"github.com/petershivachi/go_transact/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomUser(t *testing.T) User {
	args := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomCurrency(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.HashedPassword, user.HashedPassword)
	require.Equal(t, args.Username, user.FullName)
	require.Equal(t, args.Email, user.Email)

	require.NotZero(t, user.Username)
	require.NotZero(t, user.CreatedAt)

	return user
}
