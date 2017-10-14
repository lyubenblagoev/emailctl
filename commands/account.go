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

	BuildCommand(c, listAccounts, "list <domain_name>", "List accounts for specific domain", ArgsOption(1), AliasOption("l"))
	BuildCommand(c, showAccount, "show <domain_name> <account_name>", "Show information about specific account", ArgsOption(2), AliasOption("s"))
	BuildCommand(c, addAccount, "add <domain_name> <account_name>", "Add new account", ArgsOption(2), AliasOption("a"))
	BuildCommand(c, deleteAccount, "delete <domain_name> <account_name>", "Delete the account specified by domain_name and account_name", ArgsOption(2), AliasOption("rm"))
	BuildCommand(c, disableAccount, "disable <domain_name> <account_name>", "Disable the account specified by domain_name and account_name", ArgsOption(2), AliasOption("d"))
	BuildCommand(c, enableAccount, "enable <domain_name> <account_name>", "Enable the account specified by domain_name and account_name", ArgsOption(2), AliasOption("e"))
	BuildCommand(c, changeAccountPassword, "password <domain_name> <account_name>", "Change the password for the account specified by domain_name and account_name", ArgsOption(2), AliasOption("p"))

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
	password, err := emailctl.ReadPassword()
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
	password, err := emailctl.ReadPassword()
	if err != nil {
		return err
	}
	return client.Accounts.ChangePassword(domain, username, password)
}
