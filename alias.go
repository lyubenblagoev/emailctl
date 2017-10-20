package emailctl

import (
	"fmt"
	"strings"

	"github.com/lyubenblagoev/goprsc"
)

// Alias is a wrapper for goprsc.Alias.
type Alias struct {
	*goprsc.Alias
}

// AliasService is an interface for interacting with the PostfixRestServer alias API.
type AliasService interface {
	List(domain string) ([]Alias, error)
	Get(domain, alias string) ([]Alias, error)
	GetForEmail(domain, alias, email string) (*Alias, error)
	Create(domain, alias, email string) error
	Disable(domain, alias, email string) error
	Enable(domain, alias, email string) error
	Delete(domain, alias, email string) error
	Rename(domain, alias, email, newName string) error
	RenameAll(domain, alias, newName string) error
}

type aliasService struct {
	client *goprsc.Client
}

// NewAliasService builds an AliasService instance.
func NewAliasService(client *goprsc.Client) AliasService {
	return &aliasService{
		client: client,
	}
}

func (s *aliasService) List(domain string) ([]Alias, error) {
	aliases, err := s.client.Aliases.List(domain)
	if err != nil {
		return nil, err
	}

	list := make([]Alias, len(aliases))
	for i := range aliases {
		list[i] = Alias{Alias: &aliases[i]}
	}
	return list, nil
}

func (s *aliasService) Get(domain, alias string) ([]Alias, error) {
	if err := ValidateEmailAddress(alias, domain); err != nil {
		return nil, err
	}

	aliases, err := s.client.Aliases.Get(domain, alias)
	if err != nil {
		return nil, err
	}

	list := make([]Alias, len(aliases))
	for i := range aliases {
		list[i] = Alias{Alias: &aliases[i]}
	}
	return list, nil
}

func (s *aliasService) GetForEmail(domain, alias, email string) (*Alias, error) {
	if err := ValidateEmailAddress(alias, domain); err != nil {
		return nil, err
	}
	parts := strings.Split(email, "@")
	if err := ValidateEmailAddress(parts[0], parts[1]); err != nil {
		return nil, err
	}

	a, err := s.client.Aliases.GetForEmail(domain, alias, email)
	if err != nil {
		return nil, err
	}

	return &Alias{Alias: a}, nil
}

func (s *aliasService) Create(domain, alias, email string) error {
	return s.client.Aliases.Create(domain, alias, email)
}

func (s *aliasService) Delete(domain, alias, email string) error {
	return s.client.Aliases.Delete(domain, alias, email)
}

func (s *aliasService) Enable(domain, alias, email string) error {
	return s.setEnabled(domain, alias, email, true)
}

func (s *aliasService) Disable(domain, alias, email string) error {
	return s.setEnabled(domain, alias, email, false)
}

func (s *aliasService) setEnabled(domain, alias, email string, enabled bool) error {
	if err := ValidateEmailAddress(alias, domain); err != nil {
		return fmt.Errorf("Invalid alias email: %s: %v", fmt.Sprintf("%s@%s", alias, domain), err)
	}
	parts := strings.Split(email, "@")
	if err := ValidateEmailAddress(parts[0], parts[1]); err != nil {
		return fmt.Errorf("Invalid recipient address: %s: %v", email, err)
	}

	a, err := s.client.Aliases.GetForEmail(domain, alias, email)
	if err != nil {
		return err
	}

	ur := &goprsc.AliasUpdateRequest{
		Name:    alias,
		Email:   a.Email,
		Enabled: enabled,
	}
	return s.client.Aliases.Update(domain, alias, email, ur)
}

func (s *aliasService) Rename(domain, alias, email, newName string) error {
	if err := ValidateEmailAddress(newName, domain); err != nil {
		return err
	}

	ur := &goprsc.AliasUpdateRequest{
		Name:  newName,
		Email: email,
	}
	return s.client.Aliases.Update(domain, alias, email, ur)
}

func (s *aliasService) RenameAll(domain, alias, newName string) error {
	if err := ValidateEmailAddress(newName, domain); err != nil {
		return err
	}

	aliases, err := s.client.Aliases.Get(domain, alias)
	if err != nil {
		return err
	}
	for _, a := range aliases {
		ur := &goprsc.AliasUpdateRequest{
			Name:  newName,
			Email: a.Email,
		}
		if err := s.client.Aliases.Update(domain, alias, a.Email, ur); err != nil {
			return err
		}
	}

	return nil
}
