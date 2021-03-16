package emailctl

import (
	"fmt"

	"github.com/lyubenblagoev/goprsc"
	"github.com/spf13/viper"
)

// Client is the Postfix REST server client.
type Client struct {
	client *goprsc.Client

	Auth       *AuthService
	Domains    *DomainService
	Accounts   *AccountService
	Aliases    *AliasService
	InputBccs  *InputBccService
	OutputBccs *OutputBccService
}

// GetLogin returns the user login associated with the client
func (c *Client) GetLogin() string {
	return c.client.Login
}

// GetAuthToken returns the authentication token associated with the client
func (c *Client) GetAuthToken() string {
	return c.client.AuthToken
}

// GetRefreshToken returns the refresh token associated with the client
func (c *Client) GetRefreshToken() string {
	return c.client.RefreshToken
}

type service struct {
	client *goprsc.Client
}

// NewClient creates an instance of Client.
func NewClient() (*Client, error) {
	goprscClient, err := newGoprscClient()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Postfix REST Server API client: %s", err)
	}

	c := &Client{client: goprscClient}
	s := service{client: goprscClient} // Reuse the same structure instead of allocating one for each service
	c.Auth = (*AuthService)(&s)
	c.Domains = (*DomainService)(&s)
	c.Accounts = (*AccountService)(&s)
	c.Aliases = (*AliasService)(&s)
	// Allocate separate structs for the BCC services as they have different state
	c.InputBccs = NewInputBccService(goprscClient)
	c.OutputBccs = NewOutputBccService(goprscClient)
	return c, nil
}

func newGoprscClient() (*goprsc.Client, error) {
	host := viper.GetString("host")
	port := viper.GetString("port")
	useHTTPS := viper.GetBool("https")

	login := viper.GetString("login")
	authToken := viper.GetString("authToken")
	refreshToken := viper.GetString("refreshToken")

	var options []goprsc.ClientOption
	options = append(options, goprsc.UserAgentOption("emailctl"))
	options = append(options, goprsc.HostOption(host))
	options = append(options, goprsc.PortOption(port))
	if useHTTPS {
		options = append(options, goprsc.HTTPSProtocolOption())
	}
	if len(authToken) > 0 {
		options = append(options, goprsc.AuthOption(login, authToken, refreshToken))
	}

	return goprsc.NewClientWithOptions(nil, options...)
}
