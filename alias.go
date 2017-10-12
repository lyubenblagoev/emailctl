package emailctl

import "github.com/lyubenblagoev/goprsc"

// Alias is a wrapper for goprsc.Alias.
type Alias struct {
	*goprsc.Alias
}

// AliasService is an interface for interacting with the PostfixRestServer alias API.
type AliasService interface {
	List(domain string) ([]Alias, error)
	Get(domain string, alias string) (*Alias, error)
	Create(domain, alias, email string) error
	Update(domain, alias string, ur *goprsc.AliasUpdateRequest) error
	Delete(domain, alias string) error
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

func (s *aliasService) Get(domain string, alias string) (*Alias, error) {
	a, err := s.client.Aliases.Get(domain, alias)
	if err != nil {
		return nil, err
	}

	return &Alias{Alias: a}, err
}

func (s *aliasService) Create(domain, alias, email string) error {
	return s.client.Aliases.Create(domain, alias, email)
}

func (s *aliasService) Update(domain, alias string, ur *goprsc.AliasUpdateRequest) error {
	return s.client.Aliases.Update(domain, alias, ur)
}

func (s *aliasService) Delete(domain, alias string) error {
	return s.client.Aliases.Delete(domain, alias)
}
