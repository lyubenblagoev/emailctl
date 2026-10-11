package emailctl

import "github.com/lyubenblagoev/goprsc"

// DomainsPage contains one page of domains and the API's pagination metadata.
type DomainsPage struct {
	Items []Domain
	goprsc.PaginationMetadata
}

// AccountsPage contains one page of accounts and the API's pagination metadata.
type AccountsPage struct {
	Items []Account
	goprsc.PaginationMetadata
}

// AliasesPage contains one page of aliases or recipients and the API's pagination metadata.
type AliasesPage struct {
	Items []Alias
	goprsc.PaginationMetadata
}
