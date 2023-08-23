package service

import (
	"context"
	"fmt"
)

type service struct {
	storage storage
}

//go:generate mockery --structname Storage --name storage
type storage interface {
	SetBalance(ctx context.Context, b Balance) error
	GetBalance(ctx context.Context, u User) (Balance, error)
}

func New(storage storage) *service {
	return &service{storage: storage}
}

// UpBalance increases the balance of a user.
func (s *service) UpBalance(ctx context.Context, c BalanceChange) error {
	balance, err := s.storage.GetBalance(ctx, c.User)
	if err != nil && err != ErrUserNotExist {
		return fmt.Errorf("get balance from storage: %w", err)
	}

	balance.Balance += c.Changing
	if balance.Balance < 0 {
		return ErrInvalidNewBalance
	}

	if err := s.storage.SetBalance(ctx, balance); err != nil {
		return fmt.Errorf("set balance in storage: %w", err)
	}

	return nil
}

// DownBalance decreases the balance of a user.
func (s *service) DownBalance(ctx context.Context, c BalanceChange) error {
	balance, err := s.storage.GetBalance(ctx, c.User)
	if err != nil {
		return fmt.Errorf("get balance from storage: %w", err)
	}

	balance.Balance -= c.Changing
	if balance.Balance < 0 {
		return ErrInvalidNewBalance
	}

	if err := s.storage.SetBalance(ctx, balance); err != nil {
		return fmt.Errorf("set balance in storage: %w", err)
	}

	return nil
}

// GetBalance gets the balance of a user from storage.
func (s *service) GetBalance(ctx context.Context, u User) (Balance, error) {
	balance, err := s.storage.GetBalance(ctx, u)
	if err != nil {
		return Balance{}, fmt.Errorf("get balance from storage: %w", err)
	}
	return balance, nil
}
