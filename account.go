package emailctl

import (
	"fmt"
	"net/mail"

	"github.com/lyubenblagoev/goprsc"
)

// Account is a wrapper for goprsc.Account.
type Account struct {
	*goprsc.Account
}

// AccountService is an interface for interacting with the PostfixRestServer account API.
type AccountService interface {
	List(domain string) ([]Account, error)
	Get(domain string, username string) (*Account, error)
	Create(domain, username, password string) error
	Update(domain, username string, req *goprsc.AccountUpdateRequest) error
	Delete(domain, username string) error
}

type accountService struct {
	client *goprsc.Client
}

// NewAccountService builds an AccountService instance.
func NewAccountService(client *goprsc.Client) AccountService {
	return &accountService{
		client: client,
	}
}

func (s *accountService) List(domain string) ([]Account, error) {
	accounts, err := s.client.Accounts.List(domain)
	if err != nil {
		return nil, err
	}

	list := make([]Account, len(accounts))
	for i := range accounts {
		a := accounts[i]
		list[i] = Account{Account: &a}
	}

	return list, nil
}

func (s *accountService) Get(domain string, username string) (*Account, error) {
	a, err := s.client.Accounts.Get(domain, username)
	if err != nil {
		return nil, err
	}

	return &Account{Account: a}, nil
}

func (s *accountService) Create(domain, username, password string) error {
	return s.client.Accounts.Create(domain, username, password)
}

func (s *accountService) Update(domain, username string, req *goprsc.AccountUpdateRequest) error {
	return s.client.Accounts.Update(domain, username, req)
}

func (s *accountService) Delete(domain, username string) error {
	return s.client.Accounts.Delete(domain, username)
}

// ValidateEmailAddress validates the email address using the given account name and domain
// and returns an error if the name is invalid.
func ValidateEmailAddress(name, domain string) error {
	email := fmt.Sprintf("%s@%s", name, domain)
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email address: '%s'", email)
	}
	return nil
}
