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

// List retrieves one page of aliases for the specified domain.
func (s *AliasService) List(domain string, options ...goprsc.PaginationOptions) (*AliasesPage, error) {
	aliases, err := s.client.Aliases.List(domain, options...)
	if err != nil {
		return nil, err
	}

	list := make([]Alias, len(aliases.Items))
	for i := range aliases.Items {
		a := aliases.Items[i]
		list[i] = Alias{Alias: &a}
	}
	return &AliasesPage{Items: list, PaginationMetadata: aliases.PaginationMetadata}, nil
}

// Get retrieves one page of recipients for the specified alias.
func (s *AliasService) Get(domain, alias string, options ...goprsc.PaginationOptions) (*AliasesPage, error) {
	if err := ValidateEmailFromParts(alias, domain); err != nil {
		return nil, err
	}

	aliases, err := s.client.Aliases.Get(domain, alias, options...)
	if err != nil {
		return nil, err
	}

	list := make([]Alias, len(aliases.Items))
	for i := range aliases.Items {
		a := aliases.Items[i]
		list[i] = Alias{Alias: &a}
	}
	return &AliasesPage{Items: list, PaginationMetadata: aliases.PaginationMetadata}, nil
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

	aliases, err := s.listAll(domain, alias)
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

func (s *AliasService) listAll(domain, alias string) ([]goprsc.Alias, error) {
	var aliases []goprsc.Alias
	var options []goprsc.PaginationOptions
	lastPage := 0
	for {
		page, err := s.client.Aliases.Get(domain, alias, options...)
		if err != nil {
			return nil, err
		}
		if page.Page <= lastPage {
			return nil, fmt.Errorf("invalid alias pagination response: page %d followed page %d", page.Page, lastPage)
		}
		lastPage = page.Page
		aliases = append(aliases, page.Items...)
		if len(aliases) >= page.Total {
			return aliases, nil
		}
		if page.Page <= 0 || page.PageSize <= 0 || len(page.Items) == 0 {
			return nil, fmt.Errorf("invalid alias pagination response: page %d, page size %d, total %d", page.Page, page.PageSize, page.Total)
		}
		options = []goprsc.PaginationOptions{{Page: page.Page + 1, PageSize: page.PageSize}}
	}
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

	a, err := s.client.Aliases.GetForEmail(domain, alias, email)
	if err != nil {
		return err
	}
	ur := &goprsc.AliasUpdateRequest{
		Name:    newName,
		Email:   email,
		Enabled: a.Enabled,
	}
	return s.client.Aliases.Update(domain, alias, email, ur)
}

// RenameAll renames the username part of the specified aliases (for all recipients attached to the alias).
func (s *AliasService) RenameAll(domain, alias, newName string) error {
	if err := ValidateEmailFromParts(newName, domain); err != nil {
		return err
	}

	aliases, err := s.listAll(domain, alias)
	if err != nil {
		return err
	}
	for _, a := range aliases {
		ur := &goprsc.AliasUpdateRequest{
			Name:    newName,
			Email:   a.Email,
			Enabled: a.Enabled,
		}
		if err := s.client.Aliases.Update(domain, alias, a.Email, ur); err != nil {
			return err
		}
	}

	return nil
}
