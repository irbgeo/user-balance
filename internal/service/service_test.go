package service_test

import (
	"context"
	"testing"

	"github.com/irbgeo/user-balance/internal/service"
	"github.com/irbgeo/user-balance/internal/service/mocks"
	"github.com/stretchr/testify/require"
)

// TODO: занести переменные внутрь тест кейсов для более гибкойго управления
var (
	testUser = service.User{
		ID: 1,
	}

	testBalance = service.Balance{
		User:    testUser,
		Balance: 2,
	}

	testChange = service.BalanceChange{
		User:     testUser,
		Changing: 1,
	}

	testUpperBalance = service.Balance{
		User:    testUser,
		Balance: testBalance.Balance + testChange.Changing,
	}

	testDownedBalance = service.Balance{
		User:    testUser,
		Balance: testBalance.Balance - testChange.Changing,
	}

	testChangeBig = service.BalanceChange{
		User:     testUser,
		Changing: 10,
	}

	expectedError = service.ErrInvalidNewBalance
)

func TestUpBalance(t *testing.T) {
	storage := &mocks.Storage{}
	storage.On("GetBalance", context.Background(), testUser).Return(testBalance, nil)
	storage.On("SetBalance", context.Background(), testUpperBalance).Return(nil)

	svc := service.New(storage)

	err := svc.UpBalance(context.Background(), testChange)
	require.NoError(t, err)
}

func TestDownBalanceSuccess(t *testing.T) {
	storage := &mocks.Storage{}
	storage.On("GetBalance", context.Background(), testUser).Return(testBalance, nil)
	storage.On("SetBalance", context.Background(), testDownedBalance).Return(nil)

	svc := service.New(storage)

	err := svc.DownBalance(context.Background(), testChange)
	require.NoError(t, err)
}

func TestDownBalanceFail(t *testing.T) {
	storage := &mocks.Storage{}
	storage.On("GetBalance", context.Background(), testUser).Return(testBalance, nil)

	svc := service.New(storage)

	err := svc.DownBalance(context.Background(), testChangeBig)
	require.Equal(t, expectedError, err)
}

func TestGetBalance(t *testing.T) {
	storage := &mocks.Storage{}
	storage.On("GetBalance", context.Background(), testUser).Return(testBalance, nil)

	svc := service.New(storage)

	balance, err := svc.GetBalance(context.Background(), testUser)
	require.NoError(t, err)
	require.Equal(t, testBalance, balance)
}
