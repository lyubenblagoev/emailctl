package emailctl

import "github.com/lyubenblagoev/goprsc"

// Domain is a wrapper for goprsc.Domain
type Domain struct {
	*goprsc.Domain
}

// DomainService is an interface for interacting with the PostfixRestServer domain API.
type DomainService interface {
	List() ([]Domain, error)
	Get(string) (*Domain, error)
	Create(string) error
	Update(string, *goprsc.DomainUpdateRequest) error
	Delete(string) error
}

type domainService struct {
	client *goprsc.Client
}

// NewDomainService builds a DomainService instance.
func NewDomainService(client *goprsc.Client) DomainService {
	return &domainService{
		client: client,
	}
}

func (ds *domainService) List() ([]Domain, error) {
	domains, err := ds.client.Domains.List()
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

func (ds *domainService) Get(name string) (*Domain, error) {
	d, err := ds.client.Domains.Get(name)
	if err != nil {
		return nil, err
	}
	return &Domain{Domain: d}, nil
}

func (ds *domainService) Create(name string) error {
	return ds.client.Domains.Create(name)
}

func (ds *domainService) Update(name string, req *goprsc.DomainUpdateRequest) error {
	return ds.client.Domains.Update(name, req)
}

func (ds *domainService) Delete(name string) error {
	return ds.client.Domains.Delete(name)
}
