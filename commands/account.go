package commands

import (
	"fmt"

	"github.com/lyubenblagoev/emailctl"
	"github.com/spf13/cobra"
)

// CreateAccountCommand creates an account command with all its sub-commands.
func CreateAccountCommand() *Command {
	c := &Command{
		Command: &cobra.Command{
			Use:   "account",
			Short: "Account commands",
			Long:  "Account is used to access account commands",
		},
	}

	BuildCommand(c, listAccounts, "list <domain-name>", "List all accounts", ArgsOption(1), AliasOption("l"))
	BuildCommand(c, showAccount, "show <domain-name> <name>", "Show specific account", ArgsOption(2), AliasOption("s"))
	BuildCommand(c, addAccount, "add <domain-name> <name>", "Add a new account", ArgsOption(2), AliasOption("a"))
	BuildCommand(c, deleteAccount, "delete <domain-name> <name>", "Delete an account", ArgsOption(2), AliasOption("rm"))
	BuildCommand(c, disableAccount, "disable <domain-name> <name>", "Disable an account", ArgsOption(2), AliasOption("d"))
	BuildCommand(c, enableAccount, "enable <domain-name> <name>", "Enable an account", ArgsOption(2), AliasOption("e"))
	BuildCommand(c, renameAccount, "rename <domain-name> <name> <new_name>", "Rename account", ArgsOption(3), AliasOption("r"))
	BuildCommand(c, changeAccountPassword, "password <domain-name> <name>", "Change account password", ArgsOption(2), AliasOption("p"))

	return c
}

func listAccounts(client *emailctl.Client, args []string) error {
	domain := args[0]
	accounts, err := client.Accounts.List(domain)
	if err != nil {
		return err
	}

	fmt.Printf("Accounts for '%s':\n", domain)
	fmt.Printf("%-5s%-30s%-10s%-12s%-12s\n", "ID", "Email Address", "Enabled", "Created", "Updated")
	for _, a := range accounts {
		email := fmt.Sprintf("%s@%s", a.Username, domain)
		fmt.Printf("%-5d%-30s%-10t%-12s%-12s\n", a.ID, email, a.Enabled, a.Created.Format("2006-01-02"), a.Updated.Format("2006-01-02"))
	}

	return nil
}

func showAccount(client *emailctl.Client, args []string) error {
	domain, username := args[0], args[1]
	account, err := client.Accounts.Get(domain, username)
	if err != nil {
		return err
	}

	fmt.Printf("%-12s:%30d\n", "ID", account.ID)
	fmt.Printf("%-12s:%30s\n", "Email", fmt.Sprintf("%s@%s", username, domain))
	fmt.Printf("%-12s:%30t\n", "Enabled", account.Enabled)
	fmt.Printf("%-12s:%30s\n", "Created", account.Created.Format("2006-01-02"))
	fmt.Printf("%-12s:%30s\n", "Updated", account.Updated.Format("2006-01-02"))

	return nil
}

func addAccount(client *emailctl.Client, args []string) error {
	domain, username := args[0], args[1]
	password, err := emailctl.ReadAndConfirmPassword()
	if err != nil {
		return err
	}
	return client.Accounts.Create(domain, username, password)
}

func deleteAccount(client *emailctl.Client, args []string) error {
	domain, username := args[0], args[1]
	return client.Accounts.Delete(domain, username)
}

func enableAccount(client *emailctl.Client, args []string) error {
	domain, username := args[0], args[1]
	return client.Accounts.Enable(domain, username)
}

func disableAccount(client *emailctl.Client, args []string) error {
	domain, username := args[0], args[1]
	return client.Accounts.Disable(domain, username)
}

func renameAccount(client *emailctl.Client, args []string) error {
	domain, username, newName := args[0], args[1], args[2]
	return client.Accounts.Rename(domain, username, newName)
}

func changeAccountPassword(client *emailctl.Client, args []string) error {
	domain, username := args[0], args[1]
	password, err := emailctl.ReadAndConfirmPassword()
	if err != nil {
		return err
	}
	return client.Accounts.ChangePassword(domain, username, password)
}
