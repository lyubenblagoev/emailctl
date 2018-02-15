package emailctl

import "github.com/lyubenblagoev/goprsc"

// Account is a wrapper for goprsc.Account.
type Account struct {
	*goprsc.Account
}

// AccountService is an interface for interacting with the PostfixRestServer account API.
type AccountService interface {
	List(domain string) ([]Account, error)
	Get(domain string, username string) (*Account, error)
	Create(domain, username, password string) error
	Delete(domain, username string) error
	Enable(domain, username string) error
	Disable(domain, username string) error
	Rename(domain, old, new string) error
	ChangePassword(domain, username, password string) error
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
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return nil, err
	}

	a, err := s.client.Accounts.Get(domain, username)
	if err != nil {
		return nil, err
	}

	return &Account{Account: a}, nil
}

func (s *accountService) Create(domain, username, password string) error {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return err
	}

	return s.client.Accounts.Create(domain, username, password)
}

func (s *accountService) Delete(domain, username string) error {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return err
	}

	return s.client.Accounts.Delete(domain, username)
}

func (s *accountService) Enable(domain, username string) error {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return err
	}

	ur := &goprsc.AccountUpdateRequest{
		Username: username,
		Enabled:  true,
	}
	return s.client.Accounts.Update(domain, username, ur)
}

func (s *accountService) Disable(domain, username string) error {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return err
	}

	ur := &goprsc.AccountUpdateRequest{
		Username: username,
		Enabled:  false,
	}
	return s.client.Accounts.Update(domain, username, ur)
}

func (s *accountService) Rename(domain, old, new string) error {
	usernames := []string{old, new}
	for _, u := range usernames {
		if err := ValidateEmailFromParts(u, domain); err != nil {
			return err
		}
	}

	ur := &goprsc.AccountUpdateRequest{
		Username: new,
	}
	return s.client.Accounts.Update(domain, old, ur)
}

func (s *accountService) ChangePassword(domain, username, password string) error {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return err
	}

	ur := &goprsc.AccountUpdateRequest{
		Username:        username,
		Password:        password,
		ConfirmPassword: password,
	}
	return s.client.Accounts.Update(domain, username, ur)
}
