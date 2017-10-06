package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CreateVersionCommand creates a version command
func CreateVersionCommand() *Command {
	c := &Command{
		Command: &cobra.Command{
			Use:   "version",
			Short: "Prints the version number of emailctl",
			Long:  `Prints the version number of emailctl.`,

			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(EmailctlVersion.FullVersion())
			},
		},
	}
	return c
}
