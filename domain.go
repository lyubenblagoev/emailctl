package emailctl

import "github.com/lyubenblagoev/goprsc"

// Domain is a wrapper for goprsc.Domain
type Domain struct {
	*goprsc.Domain
}

// DomainService handles communication with the domain API of the Postfix REST Server.
type DomainService service

// List retrieves all domains.
func (s *DomainService) List() ([]Domain, error) {
	domains, err := s.client.Domains.List()
	if err != nil {
		return nil, err
	}

	list := make([]Domain, len(domains))
	for i := range domains {
		d := domains[i]
		list[i] = Domain{Domain: &d}
	}

	return list, nil
}

// Get retrieves a domain with the specified domain name
func (s *DomainService) Get(name string) (*Domain, error) {
	d, err := s.client.Domains.Get(name)
	if err != nil {
		return nil, err
	}
	return &Domain{Domain: d}, nil
}

// Create creates a new domain with the specified domain name.
func (s *DomainService) Create(name string) error {
	return s.client.Domains.Create(name)
}

// Delete deletes the domain with the specified name.
func (s *DomainService) Delete(name string) error {
	return s.client.Domains.Delete(name)
}

// Rename renames domain with domain name 'old' to 'new'.
func (s *DomainService) Rename(old, new string) error {
	ur := &goprsc.DomainUpdateRequest{Name: new}
	return s.client.Domains.Update(old, ur)
}

// Enable enables the domain specified by 'name'.
func (s *DomainService) Enable(name string) error {
	ur := &goprsc.DomainUpdateRequest{Name: name, Enabled: true}
	return s.client.Domains.Update(name, ur)
}

// Disable disables the domain specified by 'name'.
func (s *DomainService) Disable(name string) error {
	ur := &goprsc.DomainUpdateRequest{Name: name, Enabled: false}
	return s.client.Domains.Update(name, ur)
}
