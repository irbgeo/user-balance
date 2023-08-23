package storage

import (
	"context"
	"strconv"

	"github.com/irbgeo/user-balance/internal/service"
)

type storage struct {
	driver driver
}

type driver interface {
	Set(ctx context.Context, key, value string) error
	Get(ctx context.Context, key string) (string, error)
}

func New(driver driver) *storage {
	return &storage{driver: driver}
}

// SetBalance updates the balance of a user in the storage.
func (s *storage) SetBalance(ctx context.Context, b service.Balance) error {
	key := getKey(b.User)
	value := strconv.Itoa(b.Balance)
	return s.driver.Set(ctx, key, value)
}

// GetBalance retrieves the balance of a user from the storage.
func (s *storage) GetBalance(ctx context.Context, u service.User) (service.Balance, error) {
	key := getKey(u)
	value, err := s.driver.Get(ctx, key)
	if err != nil {
		return service.Balance{}, err
	}

	if value == "" {
		return service.Balance{User: u}, service.ErrUserNotExist
	}

	balance, err := strconv.Atoi(value)
	if err != nil {
		return service.Balance{}, err
	}

	return service.Balance{User: u, Balance: balance}, nil
}

func getKey(u service.User) string {
	return strconv.Itoa(u.ID)
}
