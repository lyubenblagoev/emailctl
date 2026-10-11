package commands

import (
	"fmt"

	"github.com/lyubenblagoev/emailctl"
	"github.com/lyubenblagoev/goprsc"
	"github.com/spf13/cobra"
)

type paginatedCommandRunner func(*emailctl.Client, []string, []goprsc.PaginationOptions) error

type paginationFlags struct {
	page     int
	pageSize int
}

func buildPaginatedCommand(parent *Command, runner paginatedCommandRunner, usage, description string, options ...CommandOption) *Command {
	pagination := &paginationFlags{}
	c := BuildCommand(parent, func(client *emailctl.Client, args []string) error {
		return runner(client, args, pagination.options())
	}, usage, description, options...)
	pagination.addFlags(c)
	return c
}

func (p *paginationFlags) addFlags(c *Command) {
	c.Flags().IntVar(&p.page, "page", 0, "One-based page number (default: 1)")
	c.Flags().IntVar(&p.pageSize, "page-size", 0, "Positive results per page (default: server-configured; maximum: server-configured)")

	previousPreRun := c.PreRunE
	c.PreRunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("page") && p.page < 1 {
			return fmt.Errorf("--page must be a positive integer")
		}
		if cmd.Flags().Changed("page-size") && p.pageSize < 1 {
			return fmt.Errorf("--page-size must be a positive integer")
		}
		if previousPreRun != nil {
			return previousPreRun(cmd, args)
		}
		return nil
	}
}

func (p *paginationFlags) options() []goprsc.PaginationOptions {
	if p.page == 0 && p.pageSize == 0 {
		return nil
	}
	return []goprsc.PaginationOptions{{Page: p.page, PageSize: p.pageSize}}
}
