package commands

import (
	"fmt"

	"github.com/lyubenblagoev/emailctl"
	"github.com/spf13/cobra"
)

type bccType uint8

const (
	sender = iota
	recipient
)

// CreateSenderBccCommand creates the sender-bcc command and its sub-commands.
func CreateSenderBccCommand() *Command {
	c := &Command{
		Command: &cobra.Command{
			Use:   "sender-bcc",
			Short: "sender-bcc commands",
			Long:  "sender-bcc is used to access sender BCC commands",
		},
	}

	BuildCommand(c, showBcc(sender), "show <domain-name> <name>", "Show BCC for specific account", ArgsOption(2), AliasOption("s"))
	BuildCommand(c, addBcc(sender), "add <domain-name> <name> <recipient-email>", "Set BCC for specified account", ArgsOption(3), AliasOption("a"))
	BuildCommand(c, deleteBcc(sender), "delete <domain-name> <name>", "Delete the sender BCC for specified account", ArgsOption(2), AliasOption("a"))
	BuildCommand(c, enableBcc(sender), "enable <domain-name> <name>", "Enable the sender BCC for specified account", ArgsOption(2), AliasOption("a"))
	BuildCommand(c, disableBcc(sender), "disable <domain-name> <name>", "Disable the sender BCC for specified account", ArgsOption(2), AliasOption("a"))
	BuildCommand(c, changeBccRecipient(sender), "change-recipient <domain-name> <name> <recipient-email>", "Change the recipient email address for existing BCC", ArgsOption(3), AliasOption("a"))

	return c
}

// CreateRecipientBccCommand creates the recipient-bcc command and its sub-commands.
func CreateRecipientBccCommand() *Command {
	c := &Command{
		Command: &cobra.Command{
			Use:   "recipient-bcc",
			Short: "recipient-bcc commands",
			Long:  "recipient-bcc is used to access recipient BCC commands",
		},
	}

	BuildCommand(c, showBcc(recipient), "show <domain-name> <name>", "Show BCC for specific account", ArgsOption(2), AliasOption("s"))
	BuildCommand(c, addBcc(recipient), "add <domain-name> <name> <recipient-email>", "Add BCC for the specified account", ArgsOption(3), AliasOption("a"))
	BuildCommand(c, deleteBcc(recipient), "delete <domain-name> <name>", "Delete the recipient BCC for specified account", ArgsOption(2), AliasOption("a"))
	BuildCommand(c, enableBcc(recipient), "enable <domain-name> <name>", "Enable the sender BCC for specified account", ArgsOption(2), AliasOption("a"))
	BuildCommand(c, disableBcc(recipient), "disable <domain-name> <name>", "Disable the sender BCC for specified account", ArgsOption(2), AliasOption("a"))
	BuildCommand(c, changeBccRecipient(recipient), "change-recipient <domain-name> <name> <recipient-email>", "Change the recipient email address for existing BCC", ArgsOption(3), AliasOption("a"))

	return c
}

func showBcc(typ bccType) CommandRunner {
	return func(client *emailctl.Client, args []string) error {
		domain, username := args[0], args[1]
		bcc, err := getService(client, typ).Get(domain, username)
		if err != nil {
			return err
		}

		fmt.Printf("%-12s:%30d\n", "ID", bcc.ID)
		fmt.Printf("%-12s:%30s\n", "Email", bcc.Email)
		fmt.Printf("%-12s:%30t\n", "Enabled", bcc.Enabled)
		fmt.Printf("%-12s:%30s\n", "Created", bcc.Created.Format("2006-01-02"))
		fmt.Printf("%-12s:%30s\n", "Updated", bcc.Updated.Format("2006-01-02"))

		return nil
	}
}

func addBcc(typ bccType) CommandRunner {
	return func(client *emailctl.Client, args []string) error {
		domain, username, email := args[0], args[1], args[2]
		return getService(client, typ).Create(domain, username, email)
	}
}

func deleteBcc(typ bccType) CommandRunner {
	return func(client *emailctl.Client, args []string) error {
		domain, username := args[0], args[1]
		return getService(client, typ).Delete(domain, username)
	}
}

func enableBcc(typ bccType) CommandRunner {
	return func(client *emailctl.Client, args []string) error {
		domain, username := args[0], args[1]
		return getService(client, typ).Enable(domain, username)
	}
}

func disableBcc(typ bccType) CommandRunner {
	return func(client *emailctl.Client, args []string) error {
		domain, username := args[0], args[1]
		return getService(client, typ).Disable(domain, username)
	}
}

func changeBccRecipient(typ bccType) CommandRunner {
	return func(client *emailctl.Client, args []string) error {
		domain, username, email := args[0], args[1], args[2]
		return getService(client, typ).ChangeRecipient(domain, username, email)
	}
}

func getService(client *emailctl.Client, typ bccType) emailctl.BccService {
	if typ == sender {
		return client.OutputBccs
	}
	return client.InputBccs
}
