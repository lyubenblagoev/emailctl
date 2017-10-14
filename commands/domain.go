package commands

import (
	"fmt"

	"github.com/lyubenblagoev/emailctl"
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

	BuildCommand(c, listDomains, "list", "List all domains", AliasOption("l"))
	BuildCommand(c, showDomain, "show <domain_name>", "Show information about specific domain", ArgsOption(1), AliasOption("s"))
	BuildCommand(c, addDomain, "add <domain_name>", "Add a new domain", ArgsOption(1), AliasOption("a"))
	BuildCommand(c, deleteDomain, "delete <domain_name>", "Delete the domain specified by domain_name", ArgsOption(1), AliasOption("rm"))
	BuildCommand(c, renameDomain, "rename <current_name> <new_name>", "Rename the domain specified by current_name to new_name", ArgsOption(2), AliasOption("r"))
	BuildCommand(c, disableDomain, "disable <domain_name>", "Disable the domain specified by domain_name", ArgsOption(1), AliasOption("d"))
	BuildCommand(c, enableDomain, "enable <domain_name>", "Enable the domain specified by domain_name", ArgsOption(1), AliasOption("e"))

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
	return client.Domains.Rename(oldName, newName)
}

func disableDomain(client *emailctl.Client, args []string) error {
	domainName := args[0]
	return client.Domains.Disable(domainName)
}

func enableDomain(client *emailctl.Client, args []string) error {
	domainName := args[0]
	return client.Domains.Enable(domainName)
}
