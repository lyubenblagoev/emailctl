package emailctl

import (
	"fmt"

	"github.com/lyubenblagoev/goprsc"
	"github.com/spf13/viper"
)

// Client is the Postfix REST server client.
type Client struct {
	client *goprsc.Client

	Domains    *DomainService
	Accounts   *AccountService
	Aliases    *AliasService
	InputBccs  *InputBccService
	OutputBccs *OutputBccService
}

type service struct {
	client *goprsc.Client
}

// NewClient creates an instance of Client.
func NewClient() (*Client, error) {
	goprscClient, err := getGoprscClient()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Postfix REST Server API client: %s", err)
	}

	c := &Client{client: goprscClient}
	s := service{client: goprscClient} // Reuse the same structure instead of allocating one for each service
	c.Domains = (*DomainService)(&s)
	c.Accounts = (*AccountService)(&s)
	c.Aliases = (*AliasService)(&s)
	// Allocate separate structs for the BCC services as they have different state
	c.InputBccs = NewInputBccService(goprscClient)
	c.OutputBccs = NewOutputBccService(goprscClient)
	return c, nil
}

func getGoprscClient() (*goprsc.Client, error) {
	host := viper.GetString("host")
	port := viper.GetString("port")
	useHTTPS := viper.GetBool("https")

	var options []goprsc.ClientOption
	options = append(options, goprsc.HostOption(host))
	options = append(options, goprsc.PortOption(port))
	if useHTTPS {
		options = append(options, goprsc.HTTPSProtocolOption())
	}

	return goprsc.NewClientWithOptions(nil, options...)
}
