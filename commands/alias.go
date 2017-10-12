package commands

import (
	"fmt"

	"github.com/lyubenblagoev/emailctl"
	"github.com/lyubenblagoev/goprsc"
	"github.com/spf13/cobra"
)

// CreateAliasCommand creates an alias command with all its sub-commands.
func CreateAliasCommand() *Command {
	c := &Command{
		Command: &cobra.Command{
			Use:   "alias",
			Short: "Alias commands",
			Long:  "Alias is used to access alias commands",
		},
	}

	BuildCommand(c, listAliases, "list <domain_name>", "List aliases for specific domain", ArgsOption(1), AliasOption("l"))
	BuildCommand(c, showAlias, "show <domain_name> <alias>", "Show information about specific alias", ArgsOption(2), AliasOption("s"))
	BuildCommand(c, addAlias, "add <domain_name> <alias> <email>", "Add new alias to an email", ArgsOption(3), AliasOption("a"))
	BuildCommand(c, deleteAlias, "delete <domain_name> <alias>", "Delete the alias specified by domain_name and alias", ArgsOption(2), AliasOption("rm"))
	BuildCommand(c, disableAlias, "disable <domain_name> <alias>", "Disable the alias specified by domain_name and alias", ArgsOption(2), AliasOption("d"))
	BuildCommand(c, enableAlias, "enable <domain_name> <alias>", "Enable the alias specified by domain_name and alias", ArgsOption(2), AliasOption("e"))

	return c
}

func listAliases(client *emailctl.Client, args []string) error {
	domain := args[0]
	aliases, err := client.Aliases.List(domain)
	if err != nil {
		return err
	}

	fmt.Printf("Aliases for '%s':\n", domain)
	fmt.Printf("%-5s%-30s%-30s%-10s%-12s%-12s\n", "ID", "Alias", "Email Address", "Enabled", "Created", "Updated")
	for _, a := range aliases {
		fmt.Printf("%-5d%-30s%-30s%-10t%-12s%-12s\n", a.ID, a.Name, a.Email, a.Enabled, a.Created.Format("2006-01-02"), a.Updated.Format("2006-01-02"))
	}

	return nil
}

func showAlias(client *emailctl.Client, args []string) error {
	domain, alias := args[0], args[1]
	validateEmailAddress(domain, alias)
	a, err := client.Aliases.Get(domain, alias)
	if err != nil {
		return err
	}

	fmt.Printf("%-12s:%30d\n", "ID", a.ID)
	fmt.Printf("%-12s:%30s\n", "Alias", a.Name)
	fmt.Printf("%-12s:%30s\n", "Email", a.Email)
	fmt.Printf("%-12s:%30t\n", "Enabled", a.Enabled)
	fmt.Printf("%-12s:%30s\n", "Created", a.Created.Format("2006-01-02"))
	fmt.Printf("%-12s:%30s\n", "Updated", a.Updated.Format("2006-01-02"))

	return nil
}

func addAlias(client *emailctl.Client, args []string) error {
	domain, alias, email := args[0], args[1], args[2]
	validateEmailAddress(domain, alias)
	return client.Aliases.Create(domain, alias, email)
}

func deleteAlias(client *emailctl.Client, args []string) error {
	domain, alias := args[0], args[1]
	validateEmailAddress(domain, alias)
	return client.Aliases.Delete(domain, alias)
}

func enableAlias(client *emailctl.Client, args []string) error {
	domain, alias := args[0], args[1]
	return changeAliasEnabled(client, domain, alias, true)
}

func disableAlias(client *emailctl.Client, args []string) error {
	domain, alias := args[0], args[1]
	return changeAliasEnabled(client, domain, alias, false)
}

func changeAliasEnabled(client *emailctl.Client, domain, alias string, enabled bool) error {
	validateEmailAddress(domain, alias)
	a, err := client.Aliases.Get(domain, alias)
	if err != nil {
		return err
	}
	ur := &goprsc.AliasUpdateRequest{
		Name:    alias,
		Email:   a.Email,
		Enabled: enabled,
	}
	return client.Aliases.Update(domain, alias, ur)
}
