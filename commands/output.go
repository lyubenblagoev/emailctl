package commands

import (
	"fmt"
	"io"

	"github.com/lyubenblagoev/emailctl"
)

func writeDomainPage(w io.Writer, page *emailctl.DomainsPage) {
	fmt.Fprintf(w, "Domains (page %d, page size %d; %d total):\n", page.Page, page.PageSize, page.Total)
	fmt.Fprintf(w, "%-5s%-30s%-10s%-12s%-12s\n", "ID", "Name", "Enabled", "Created", "Updated")
	for _, domain := range page.Items {
		fmt.Fprintf(w, "%-5d%-30s%-10t%-12s%-12s\n", domain.ID, domain.Name, domain.Enabled, domain.Created.Format("2006-01-02"), domain.Updated.Format("2006-01-02"))
	}
}

func writeAccountPage(w io.Writer, name string, page *emailctl.AccountsPage) {
	fmt.Fprintf(w, "Accounts for '%s' (page %d, page size %d; %d total):\n", name, page.Page, page.PageSize, page.Total)
	fmt.Fprintf(w, "%-5s%-30s%-10s%-12s%-12s\n", "ID", "Email Address", "Enabled", "Created", "Updated")
	for _, account := range page.Items {
		email := fmt.Sprintf("%s@%s", account.Username, name)
		fmt.Fprintf(w, "%-5d%-30s%-10t%-12s%-12s\n", account.ID, email, account.Enabled, account.Created.Format("2006-01-02"), account.Updated.Format("2006-01-02"))
	}
}

func writeAliasPage(w io.Writer, title string, page *emailctl.AliasesPage) {
	fmt.Fprintf(w, "%s (page %d, page size %d; %d total):\n", title, page.Page, page.PageSize, page.Total)
	fmt.Fprintf(w, "%-5s%-30s%-30s%-10s%-12s%-12s\n", "ID", "Alias", "Email Address", "Enabled", "Created", "Updated")
	for _, alias := range page.Items {
		fmt.Fprintf(w, "%-5d%-30s%-30s%-10t%-12s%-12s\n", alias.ID, alias.Name, alias.Email, alias.Enabled, alias.Created.Format("2006-01-02"), alias.Updated.Format("2006-01-02"))
	}
}
