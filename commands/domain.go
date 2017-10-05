package commands

import (
	"fmt"

	"github.com/lyubenblagoev/emailctl"
	"github.com/lyubenblagoev/goprsc"
	"github.com/spf13/cobra"
)

// CreateDomainCommand creates a domain command with all its sub-commands.
func CreateDomainCommand() *Command {
	c := &Command{
		Command: &cobra.Command{
			Use:   "domain",
			Short: "Domain commands",
			Long:  "Domain is used to access domain commands",
		},
	}

	BuildCommand(c, listDomains, "list", "List all domains")
	BuildCommand(c, showDomain, "show <domain_name>", "Shows information about specific domain", ArgsOption(1))
	BuildCommand(c, addDomain, "add <domain_name>", "Adds a new domain", ArgsOption(1))
	BuildCommand(c, deleteDomain, "delete <domain_name>", "Deletes the domain specified by domain_name", ArgsOption(1))
	BuildCommand(c, renameDomain, "rename <current_name> <new_name>", "Renames the domain specified by current_name to new_name", ArgsOption(2))
	BuildCommand(c, disableDomain, "disable <domain_name>", "Disables the domain specified by domain_name", ArgsOption(1))
	BuildCommand(c, enableDomain, "enable <domain_name>", "Enables the domain specified by domain_name", ArgsOption(1))

	return c
}

func listDomains(client *emailctl.Client, args []string) error {
	domains, err := client.Domains.List()
	if err != nil {
		return err
	}

	fmt.Printf("Domains:\n")
	fmt.Printf("%-5s%-30s%-10s%-12s%-12s\n", "ID", "Name", "Enabled", "Created", "Updated")
	for _, d := range domains {
		fmt.Printf("%-5d%-30s%-10t%-12s%-12s\n", d.ID, d.Name, d.Enabled, d.Created.Format("2006-01-02"), d.Updated.Format("2006-01-02"))
	}

	return nil
}

func showDomain(client *emailctl.Client, args []string) error {
	name := args[0]
	d, err := client.Domains.Get(name)
	if err != nil {
		return err
	}

	fmt.Printf("%-12s:%30d\n", "ID", d.ID)
	fmt.Printf("%-12s:%30s\n", "Domain Name", name)
	fmt.Printf("%-12s:%30t\n", "Enabled", d.Enabled)
	fmt.Printf("%-12s:%30s\n", "Created", d.Created.Format("2006-01-02"))
	fmt.Printf("%-12s:%30s\n", "Updated", d.Updated.Format("2006-01-02"))

	return nil
}

func addDomain(client *emailctl.Client, args []string) error {
	name := args[0]
	return client.Domains.Create(name)
}

func deleteDomain(client *emailctl.Client, args []string) error {
	name := args[0]
	return client.Domains.Delete(name)
}

func renameDomain(client *emailctl.Client, args []string) error {
	oldName, newName := args[0], args[1]
	ur := &goprsc.DomainUpdateRequest{
		Name:    newName,
		Enabled: true,
	}
	return client.Domains.Update(oldName, ur)
}

func disableDomain(client *emailctl.Client, args []string) error {
	domainName := args[0]
	ur := &goprsc.DomainUpdateRequest{
		Name:    domainName,
		Enabled: false,
	}
	return client.Domains.Update(domainName, ur)
}

func enableDomain(client *emailctl.Client, args []string) error {
	domainName := args[0]
	ur := &goprsc.DomainUpdateRequest{
		Name:    domainName,
		Enabled: true,
	}
	return client.Domains.Update(domainName, ur)
}
