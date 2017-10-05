package commands

import (
	"github.com/lyubenblagoev/emailctl"
	"github.com/spf13/cobra"
)

// Command is a wrapper for cobra.Command.
type Command struct {
	*cobra.Command
}

// CommandRunner runs a Command using the provider client and args.
type CommandRunner func(*emailctl.Client, []string) error

// CommandOption is an option to Command.
type CommandOption func(*Command)

// BuildCommand creates a new Command.
func BuildCommand(parent *Command, runner CommandRunner, usage, description string, options ...CommandOption) *Command {
	cobraCmd := &cobra.Command{
		Use:   usage,
		Short: description,
		Long:  description,
		Run: func(cmd *cobra.Command, args []string) {
			client, err := emailctl.NewClient()
			checkErr(err)
			checkErr(runner(client, args))
		},
	}

	c := &Command{Command: cobraCmd}

	for _, o := range options {
		o(c)
	}

	if parent != nil {
		parent.AddCommand(c)
	}

	return c
}

// AddCommand add child commands to the Command.
func (c *Command) AddCommand(commands ...*Command) {
	for _, cmd := range commands {
		c.Command.AddCommand(cmd.Command)
	}
}

// ArgsOption returns a CommandOption that returns an error if there are not exactly n arguments.
func ArgsOption(n int) CommandOption {
	return func(c *Command) {
		c.Args = cobra.ExactArgs(n)
	}
}
