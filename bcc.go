package emailctl

import (
	"github.com/lyubenblagoev/goprsc"
)

// Bcc is a wrapper for goprsc.Bcc
type Bcc struct {
	*goprsc.Bcc
}

// BccService is an interface for interacting with the PostfixRestServer bcc API.
type BccService interface {
	Get(domain, username string) (*Bcc, error)
	Create(domain, username, email string) error
	Delete(domain, username string) error
	Enable(domain, username string) error
	Disable(domain, username string) error
	ChangeRecipient(domain, username, email string) error
}

type bccService struct {
	client  *goprsc.Client
	service goprsc.BccService
}

// NewInputBccService builds a new BccService instance for interacting with the sender BCC API.
func NewInputBccService(client *goprsc.Client) BccService {
	return &bccService{
		client:  client,
		service: goprsc.NewInputBccService(client),
	}
}

// NewOutputBccService builds a new BccService instance for interacting with the recipient BCC API.
func NewOutputBccService(client *goprsc.Client) BccService {
	return &bccService{
		client:  client,
		service: goprsc.NewOutputBccService(client),
	}
}

func (s *bccService) Get(domain, username string) (*Bcc, error) {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return nil, err
	}

	b, err := s.service.Get(domain, username)
	if err != nil {
		return nil, err
	}

	return &Bcc{Bcc: b}, nil
}

func (s *bccService) Create(domain, username, email string) error {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return err
	}
	if err := ValidateEmail(email); err != nil {
		return err
	}

	return s.service.Create(domain, username, email)
}

func (s *bccService) Delete(domain, username string) error {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return err
	}
	return s.service.Delete(domain, username)
}

func (s *bccService) Enable(domain, username string) error {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return err
	}

	ur := &goprsc.BccUpdateRequest{
		Enabled: true,
	}
	return s.service.Update(domain, username, ur)
}

func (s *bccService) Disable(domain, username string) error {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return err
	}

	ur := &goprsc.BccUpdateRequest{
		Enabled: false,
	}
	return s.service.Update(domain, username, ur)
}

func (s *bccService) ChangeRecipient(domain, username, email string) error {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return err
	}
	if err := ValidateEmail(email); err != nil {
		return err
	}

	ur := &goprsc.BccUpdateRequest{
		Email:   email,
		Enabled: false,
	}
	return s.service.Update(domain, username, ur)
}
