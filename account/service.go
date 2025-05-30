package account

import (
	"context"
	"errors"

	"github.com/segmentio/ksuid"
)

type AccountIdentifier struct {
	ID    *string
	Email *string
}

type Service interface {
	PostAccount(ctx context.Context, name string) (*Account, error)
	GetAccount(ctx context.Context, AccountIdentifier string) (*Account, error)
	GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type Account struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type accountService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &accountService{r}
}

func (s *accountService) PostAccount(ctx context.Context, name string, email string) (*Account, error) {
	a := &Account{
		Name:  name,
		ID:    ksuid.New().String(),
		Email: email,
	}
	if err := s.repository.PutAccount(ctx, *a); err != nil {
		return nil, err
	}
	return a, nil
}

// this part needs editing
func (s *accountService) GetAccount(ctx context.Context, identifier AccountIdentifier) (*Account, error) {
	switch {
	case identifier.ID != nil && identifier.Email != nil:
		// If both are provided, prefer ID as it's more specific
		return s.repository.GetAccountByID(ctx, *identifier.ID)
	case identifier.ID != nil:
		return s.repository.GetAccountByID(ctx, *identifier.ID)
	case identifier.Email != nil:
		return s.repository.GetAccountByEmail(ctx, *identifier.Email)
	default:
		return nil, errors.New("must provide either ID or Email")
	}
}

func (s *accountService) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.ListAccounts(ctx, skip, take)
}
