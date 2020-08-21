package emailctl

import (
	"github.com/lyubenblagoev/goprsc"
)

// Bcc is a wrapper for goprsc.Bcc
type Bcc struct {
	*goprsc.Bcc
}

// BccService handles communication with the BCC API.
type BccService interface {
	Get(domain, username string) (*Bcc, error)
	Create(domain, username, email string) error
	Delete(domain, username string) error
	Enable(domain, username string) error
	Disable(domain, username string) error
	ChangeRecipient(domain, username, email string) error
}

type bccServiceImpl struct {
	client     *goprsc.Client
	bccService goprsc.BccService
}

// InputBccService handles communication with the recipient BCC API.
type InputBccService struct {
	*bccServiceImpl
}

// NewInputBccService builds a new InputBccServer instance for interacting with the recipient BCC API.
func NewInputBccService(client *goprsc.Client) *InputBccService {
	return &InputBccService{
		bccServiceImpl: &bccServiceImpl{
			client:     client,
			bccService: goprsc.NewIncomingBccService(client),
		},
	}
}

// OutputBccService handles communication with the sender BCC API.
type OutputBccService struct {
	*bccServiceImpl
}

// NewOutputBccService builds a new OutputBccService instance for interacting with the sender BCC API.
func NewOutputBccService(client *goprsc.Client) *OutputBccService {
	return &OutputBccService{
		bccServiceImpl: &bccServiceImpl{
			client:     client,
			bccService: goprsc.NewOutgoingBccService(client),
		},
	}
}

func (s *bccServiceImpl) Get(domain, username string) (*Bcc, error) {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return nil, err
	}

	b, err := s.bccService.Get(domain, username)
	if err != nil {
		return nil, err
	}

	return &Bcc{Bcc: b}, nil
}

func (s *bccServiceImpl) Create(domain, username, email string) error {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return err
	}
	if err := ValidateEmail(email); err != nil {
		return err
	}

	return s.bccService.Create(domain, username, email)
}

func (s *bccServiceImpl) Delete(domain, username string) error {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return err
	}
	return s.bccService.Delete(domain, username)
}

func (s *bccServiceImpl) Enable(domain, username string) error {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return err
	}

	ur := &goprsc.BccUpdateRequest{
		Enabled: true,
	}
	return s.bccService.Update(domain, username, ur)
}

func (s *bccServiceImpl) Disable(domain, username string) error {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return err
	}

	ur := &goprsc.BccUpdateRequest{
		Enabled: false,
	}
	return s.bccService.Update(domain, username, ur)
}

func (s *bccServiceImpl) ChangeRecipient(domain, username, email string) error {
	if err := ValidateEmailFromParts(username, domain); err != nil {
		return err
	}
	if err := ValidateEmail(email); err != nil {
		return err
	}

	bcc, err := s.bccService.Get(domain, username)
	if err != nil {
		return err
	}
	ur := &goprsc.BccUpdateRequest{
		Email:   email,
		Enabled: bcc.Enabled,
	}
	return s.bccService.Update(domain, username, ur)
}
