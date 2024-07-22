package module

import (
	"context"
	"fmt"
	"sync"
)

type Service interface {
	RegisterUsers() error
	CreateWalletPerUsers() error
	TopupBalance(username string, amount float64) error
	Transfer(ctx context.Context, from, to, idpotency string, amount float64) error
	GetWalletByName(from string) (*Wallet, error)
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
		mapUser: map[string]*UserWallet{
			User1: {User: nil, Wallet: nil},
			User2: {User: nil, Wallet: nil},
			User3: {User: nil, Wallet: nil},
		},
	}
}

type service struct {
	repo    Repository
	mapUser map[string]*UserWallet
	rw      sync.RWMutex
}

func (s *service) RegisterUsers() error {
	for name := range s.mapUser {
		user, err := s.repo.Register(name)
		if err != nil {
			return err
		}
		s.mapUser[name].User = user
	}
	return nil
}

func (s *service) CreateWalletPerUsers() error {
	for name, data := range s.mapUser {
		if data.User == nil {
			return fmt.Errorf("user `%s` is not created", name)
		}
		wallet, err := s.repo.CreateWallet(data.User, 10000)
		if err != nil {
			return err
		}
		s.mapUser[name].Wallet = wallet
	}
	return nil
}

func (s *service) TopupBalance(username string, amount float64) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	data, ok := s.mapUser[username]
	if !ok {
		return fmt.Errorf("user `%s` is not found", username)
	}
	return s.repo.TopupBalance(data.User, data.Wallet, amount)
}

func (s *service) Transfer(ctx context.Context, from, to, idpotency string, amount float64) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	dataFrom, ok := s.mapUser[from]
	if !ok {
		return fmt.Errorf("user `%s` is not found", from)
	}
	dataTo, ok := s.mapUser[to]
	if !ok {
		return fmt.Errorf("user `%s` is not found", to)
	}

	has, err := s.repo.ValidateIdempotency(ctx, idpotency)
	if err != nil {
		return err
	}
	if has {
		return fmt.Errorf("duplicate request")
	}

	err = s.repo.Transfer(dataFrom.User, dataTo.User, dataFrom.Wallet, dataTo.Wallet, amount)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetWalletByName(from string) (*Wallet, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	dataFrom, ok := s.mapUser[from]
	if !ok {
		return nil, fmt.Errorf("user `%s` is not found", from)
	}
	result, err := s.repo.GetWalletByID(dataFrom.Wallet.ID)
	if err != nil {
		return nil, err
	}
	return result, nil
}
