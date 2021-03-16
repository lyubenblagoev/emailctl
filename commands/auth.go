package commands

import (
	"fmt"

	"github.com/lyubenblagoev/emailctl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CreateAuthCommand creates an auth command with its subcommands.
func CreateAuthCommand() *Command {
	c := &Command{
		Command: &cobra.Command{
			Use:   "auth",
			Short: "Authentication commands",
			Long:  "Auth is used to access authentication commands",
		},
	}
	BuildCommand(c, login, "login <email-address>", "Log in using the given email address", ArgsOption(1), AliasOption("l"))
	BuildCommand(c, logout, "logout", "Log in using the given email address", AliasOption("l"))
	return c
}

func login(client *emailctl.Client, args []string) error {
	login := args[0]
	password, err := emailctl.ReadPassword("Password: ")
	if err != nil {
		return err
	}
	auth, err := client.Auth.Login(login, password)
	if err != nil {
		return err
	}
	fmt.Printf("Logged in successfully.\n")
	err = SaveAuth(login, auth.AuthToken, auth.RefreshToken)
	return err
}

func logout(client *emailctl.Client, args []string) error {
	login := viper.GetString("login")
	refreshToken := viper.GetString("refreshToken")
	CleanAuth()
	return client.Auth.Logout(login, refreshToken)
}

// SaveAuth writes active authentication tokens to the configuration file
func SaveAuth(login, token, refreshToken string) error {
	if login != "" && token != "" && refreshToken != "" {
		viper.Set("login", login)
		viper.Set("authToken", token)
		viper.Set("refreshToken", refreshToken)
		return viper.WriteConfig()
	}
	return nil
}

// CleanAuth removes saved authentication tokens
func CleanAuth() {
	viper.Set("login", "")
	viper.Set("authToken", "")
	viper.Set("refreshToken", "")
	viper.WriteConfig()
}
