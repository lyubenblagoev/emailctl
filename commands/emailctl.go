package commands

import (
	"log"
	"os"

	"github.com/lyubenblagoev/emailctl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// EmailctlVersion is emailctl's version.
var EmailctlVersion = emailctl.Version{
	Major: 0,
	Minor: 2,
	Patch: 0,
}

var cfgFile string
var client *emailctl.Client

// emailctlCommand represents the base command when called without any subcommands
var emailctlCommand = &Command{
	Command: &cobra.Command{
		Use:   "emailctl",
		Short: "emailctl is a CLI for managing Postfix Rest Server",
		Long:  `emailctl is a command line interface (CLI) to the Postfix Rest Server`,
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
func Execute() {
	checkErr(emailctlCommand.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	emailctlCommand.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.emailctl.yaml)")
	initCommands()
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".emailctl")
		viper.SetConfigType("yaml")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to determine user's home directory.", err)
	}
	viper.AddConfigPath(home)
	viper.AutomaticEnv()

	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", "8080")
	viper.SetDefault("https", false)

	checkErr(viper.ReadInConfig())

	initClient()
}

func initCommands() {
	emailctlCommand.AddCommand(CreateAuthCommand())
	emailctlCommand.AddCommand(CreateDomainCommand())
	emailctlCommand.AddCommand(CreateAccountCommand())
	emailctlCommand.AddCommand(CreateAliasCommand())
	emailctlCommand.AddCommand(CreateVersionCommand())
	emailctlCommand.AddCommand(CreateSenderBccCommand())
	emailctlCommand.AddCommand(CreateRecipientBccCommand())
}

func initClient() {
	var err error
	client, err = emailctl.NewClient()
	checkErr(err)
}
