package emailctl

import "github.com/lyubenblagoev/goprsc"

// Account is a wrapper for goprsc.Account.
type Account struct {
	*goprsc.Account
}

// AccountService handles communication with the account API on the Postfix REST Server.
type AccountService service

// List retrieves all accounts for the specified domain.
func (s *AccountService) List(domain string) ([]Account, error) {
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

// Get retrieves the account for domain 'domain' with name 'username'.
func (s *AccountService) Get(domain string, username string) (*Account, error) {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return nil, err
	}

	a, err := s.client.Accounts.Get(domain, username)
	if err != nil {
		return nil, err
	}

	return &Account{Account: a}, nil
}

// Create creates a new account in the specified domain with the given username
// and password.
func (s *AccountService) Create(domain, username, password string) error {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return err
	}

	return s.client.Accounts.Create(domain, username, password)
}

// Delete deletes the specified account.
func (s *AccountService) Delete(domain, username string) error {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return err
	}

	return s.client.Accounts.Delete(domain, username)
}

// Enable enables the specified account.
func (s *AccountService) Enable(domain, username string) error {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return err
	}

	ur := &goprsc.AccountUpdateRequest{
		Username: username,
		Enabled:  true,
	}
	return s.client.Accounts.Update(domain, username, ur)
}

// Disable disables the specified account.
func (s *AccountService) Disable(domain, username string) error {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return err
	}

	ur := &goprsc.AccountUpdateRequest{
		Username: username,
		Enabled:  false,
	}
	return s.client.Accounts.Update(domain, username, ur)
}

// Rename renames the specified account username from 'old'@domain to 'new'@domain.
func (s *AccountService) Rename(domain, old, new string) error {
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

// ChangePassword changes the password for the specified account.
func (s *AccountService) ChangePassword(domain, username, password string) error {
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
