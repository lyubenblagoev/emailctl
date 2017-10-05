package commands

import (
	"fmt"

	"github.com/lyubenblagoev/emailctl"
	"github.com/lyubenblagoev/goprsc"
	"github.com/spf13/cobra"
)

var (
	domainCommand = &cobra.Command{
		Use:   "domain",
		Short: "Domain commands",
		Long:  "Domain is used to access domain commands",
	}

	domainListCommand = &cobra.Command{
		Use:   "list",
		Short: "Lists all domains",
		Long:  "Lists all domains",
		Run:   listDomains,
	}

	domainShowCommand = &cobra.Command{
		Use:   "show <domain_name>",
		Short: "Shows information about specific domain",
		Long:  "Shows information about specific domain",
		Args:  cobra.ExactArgs(1),
		Run:   showDomain,
	}

	domainAddCommand = &cobra.Command{
		Use:   "add <domain_name>",
		Short: "Adds a new domain.",
		Long:  "Adds a new domain.",
		Args:  cobra.ExactArgs(1),
		Run:   addDomain,
	}

	domainDeleteCommand = &cobra.Command{
		Use:   "delete <domain_name>",
		Short: "Deletes the domain specified by domain_name",
		Long:  "Deletes the domain specified by domain_name",
		Args:  cobra.ExactArgs(1),
		Run:   deleteDomain,
	}

	domainRenameCommand = &cobra.Command{
		Use:   "rename <current_name> <new_name>",
		Short: "Renames the domain specified by current_name to new_name",
		Long:  "Renames the domain specified by current_name to new_name",
		Args:  cobra.ExactArgs(2),
		Run:   renameDomain,
	}

	domainDisableCommand = &cobra.Command{
		Use:   "disable <domain_name>",
		Short: "Disables the domain specified by domain_name",
		Long:  "Disables the domain specified by domain_name",
		Args:  cobra.ExactArgs(1),
		Run:   disableDomain,
	}

	domainEnableCommand = &cobra.Command{
		Use:   "enable <domain_name>",
		Short: "Enables the domain specified by domain_name",
		Long:  "Enables the domain specified by domain_name",
		Args:  cobra.ExactArgs(1),
		Run:   enableDomain,
	}
)

func init() {
	domainCommand.AddCommand(domainListCommand)
	domainCommand.AddCommand(domainShowCommand)
	domainCommand.AddCommand(domainAddCommand)
	domainCommand.AddCommand(domainDeleteCommand)
	domainCommand.AddCommand(domainRenameCommand)
	domainCommand.AddCommand(domainDisableCommand)
	domainCommand.AddCommand(domainEnableCommand)

	emailctlCommand.AddCommand(domainCommand)
}

func listDomains(cmd *cobra.Command, args []string) {
	client, err := emailctl.NewClient()
	checkErr(err)

	domains, err := client.Domains.List()
	checkErr(err)

	fmt.Printf("Domains:\n")
	fmt.Printf("%-5s%-30s%-10s%-12s%-12s\n", "ID", "Name", "Enabled", "Created", "Updated")
	for _, d := range domains {
		fmt.Printf("%-5d%-30s%-10t%-12s%-12s\n", d.ID, d.Name, d.Enabled, d.Created.Format("2006-01-02"), d.Updated.Format("2006-01-02"))
	}
}

func showDomain(cmd *cobra.Command, args []string) {
	client, err := emailctl.NewClient()
	checkErr(err)

	name := args[0]
	d, err := client.Domains.Get(name)
	checkErr(err)

	fmt.Printf("%-12s:%30d\n", "ID", d.ID)
	fmt.Printf("%-12s:%30s\n", "Domain Name", name)
	fmt.Printf("%-12s:%30t\n", "Enabled", d.Enabled)
	fmt.Printf("%-12s:%30s\n", "Created", d.Created.Format("2006-01-02"))
	fmt.Printf("%-12s:%30s\n", "Updated", d.Updated.Format("2006-01-02"))
}

func addDomain(cmd *cobra.Command, args []string) {
	client, err := emailctl.NewClient()
	checkErr(err)

	name := args[0]
	checkErr(client.Domains.Create(name))
}

func deleteDomain(cmd *cobra.Command, args []string) {
	client, err := emailctl.NewClient()
	checkErr(err)

	name := args[0]
	checkErr(client.Domains.Delete(name))
}

func renameDomain(cmd *cobra.Command, args []string) {
	client, err := emailctl.NewClient()
	checkErr(err)

	oldName, newName := args[0], args[1]
	ur := &goprsc.DomainUpdateRequest{
		Name:    newName,
		Enabled: true,
	}
	checkErr(client.Domains.Update(oldName, ur))
}

func disableDomain(cmd *cobra.Command, args []string) {
	client, err := emailctl.NewClient()
	checkErr(err)

	domainName := args[0]
	ur := &goprsc.DomainUpdateRequest{
		Name:    domainName,
		Enabled: false,
	}
	checkErr(client.Domains.Update(domainName, ur))
}

func enableDomain(cmd *cobra.Command, args []string) {
	client, err := emailctl.NewClient()
	checkErr(err)

	domainName := args[0]
	ur := &goprsc.DomainUpdateRequest{
		Name:    domainName,
		Enabled: true,
	}
	checkErr(client.Domains.Update(domainName, ur))
}
