package commands

import (
	"fmt"
	"os"

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

	buildPaginatedCommand(c, listAliases, "list <domain-name> [<name>]", "List aliases", ArgsRangeOption(1, 2), AliasOption("l"))
	BuildCommand(c, showAlias, "show <domain-name> <name> <recipient-email-address>", "Show specific alias", ArgsOption(3), AliasOption("s"))
	BuildCommand(c, addAlias, "add <domain-name> <name> <recipient-email-address>", "Add a new alias", ArgsOption(3), AliasOption("a"))
	BuildCommand(c, deleteAlias, "delete <domain-name> <name> [<recipient-email-address>]", "Delete alias(es)", ArgsRangeOption(2, 3), AliasOption("rm"))
	BuildCommand(c, disableAlias, "disable <domain-name> <name> <recipient-email-address>", "Disable specific alias", ArgsOption(3), AliasOption("d"))
	BuildCommand(c, enableAlias, "enable <domain-name> <name> <recipient-email-address>", "Enable specific alias", ArgsOption(3), AliasOption("e"))
	BuildCommand(c, renameAlias, "rename <domain-name> <name> <new-name> [<recipient-email-address>]", "Rename alias(es)", ArgsRangeOption(3, 4), AliasOption("r"))

	return c
}

func listAliases(client *emailctl.Client, args []string, options []goprsc.PaginationOptions) error {
	if len(args) == 1 {
		return listForDomain(client, args, options)
	}
	return listForAlias(client, args, options)
}

func listForDomain(client *emailctl.Client, args []string, options []goprsc.PaginationOptions) error {
	domain := args[0]
	aliases, err := client.Aliases.List(domain, options...)
	if err != nil {
		return err
	}

	writeAliasPage(os.Stdout, fmt.Sprintf("Aliases for '%s'", domain), aliases)
	return nil
}

func listForAlias(client *emailctl.Client, args []string, options []goprsc.PaginationOptions) error {
	domain, alias := args[0], args[1]
	aliases, err := client.Aliases.Get(domain, alias, options...)
	if err != nil {
		return err
	}

	writeAliasPage(os.Stdout, fmt.Sprintf("Aliases for '%s@%s'", alias, domain), aliases)
	return nil
}

func showAlias(client *emailctl.Client, args []string) error {
	domain, alias, email := args[0], args[1], args[2]
	a, err := client.Aliases.GetForEmail(domain, alias, email)
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
	return client.Aliases.Create(domain, alias, email)
}

func deleteAlias(client *emailctl.Client, args []string) error {
	domain, alias := args[0], args[1]
	if len(args) == 3 {
		email := args[2]
		return client.Aliases.Delete(domain, alias, email)
	}
	return client.Aliases.DeleteAll(domain, alias)
}

func enableAlias(client *emailctl.Client, args []string) error {
	domain, alias, email := args[0], args[1], args[2]
	return client.Aliases.Enable(domain, alias, email)
}

func disableAlias(client *emailctl.Client, args []string) error {
	domain, alias, email := args[0], args[1], args[2]
	return client.Aliases.Disable(domain, alias, email)
}

func renameAlias(client *emailctl.Client, args []string) error {
	domain, alias, newName := args[0], args[1], args[2]
	if len(args) == 4 {
		email := args[3]
		return client.Aliases.Rename(domain, alias, email, newName)
	}
	return client.Aliases.RenameAll(domain, alias, newName)
}
