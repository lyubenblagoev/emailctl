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

// AliasService handles communication with the alias API of the Postfix REST Server.
type AliasService service

// List retrieves all aliases for accounts in the specified domain.
func (s *AliasService) List(domain string) ([]Alias, error) {
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

// Get retrieves all recipients for the specified alias.
func (s *AliasService) Get(domain, alias string) ([]Alias, error) {
	if err := ValidateEmailFromParts(alias, domain); err != nil {
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

// GetForEmail retreives a specific alias.
func (s *AliasService) GetForEmail(domain, alias, email string) (*Alias, error) {
	if err := ValidateEmailFromParts(alias, domain); err != nil {
		return nil, err
	}
	parts := strings.Split(email, "@")
	if err := ValidateEmailFromParts(parts[0], parts[1]); err != nil {
		return nil, err
	}

	a, err := s.client.Aliases.GetForEmail(domain, alias, email)
	if err != nil {
		return nil, err
	}

	return &Alias{Alias: a}, nil
}

// Create assignes email to the specified alias.
func (s *AliasService) Create(domain, alias, email string) error {
	return s.client.Aliases.Create(domain, alias, email)
}

// Delete deletes the specified alias.
func (s *AliasService) Delete(domain, alias, email string) error {
	if err := ValidateEmailFromParts(alias, domain); err != nil {
		return err
	}
	return s.client.Aliases.Delete(domain, alias, email)
}

// DeleteAll deletes all recipients for a specific alias.
func (s *AliasService) DeleteAll(domain, alias string) error {
	if err := ValidateEmailFromParts(alias, domain); err != nil {
		return err
	}

	aliases, err := s.client.Aliases.Get(domain, alias)
	if err != nil {
		return err
	}
	for _, a := range aliases {
		if err := s.client.Aliases.Delete(domain, alias, a.Email); err != nil {
			return err
		}
	}

	return nil
}

// Enable enables the specified alias.
func (s *AliasService) Enable(domain, alias, email string) error {
	return s.setEnabled(domain, alias, email, true)
}

// Disable disables the specified alias.
func (s *AliasService) Disable(domain, alias, email string) error {
	return s.setEnabled(domain, alias, email, false)
}

func (s *AliasService) setEnabled(domain, alias, email string, enabled bool) error {
	if err := ValidateEmailFromParts(alias, domain); err != nil {
		return fmt.Errorf("Invalid alias email: %s: %v", fmt.Sprintf("%s@%s", alias, domain), err)
	}
	parts := strings.Split(email, "@")
	if err := ValidateEmailFromParts(parts[0], parts[1]); err != nil {
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

// Rename changes the username part of the specified alias forwarding to the specified email address.
func (s *AliasService) Rename(domain, alias, email, newName string) error {
	if err := ValidateEmailFromParts(newName, domain); err != nil {
		return err
	}

	ur := &goprsc.AliasUpdateRequest{
		Name:  newName,
		Email: email,
	}
	return s.client.Aliases.Update(domain, alias, email, ur)
}

// RenameAll renames the username part of the specified aliases (for all recipients attached to the alias).
func (s *AliasService) RenameAll(domain, alias, newName string) error {
	if err := ValidateEmailFromParts(newName, domain); err != nil {
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
