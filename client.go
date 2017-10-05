package emailctl

import (
	"fmt"

	"github.com/lyubenblagoev/goprsc"
	"github.com/spf13/viper"
)

// Client is the Postfix REST server client.
type Client struct {
	client *goprsc.Client

	Domains DomainService
}

// NewClient creates an instance of Client.
func NewClient() (*Client, error) {
	goprscClient, err := getGoprscClient()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Postfix REST Server API client: %s", err)
	}

	client := &Client{
		client:  goprscClient,
		Domains: NewDomainService(goprscClient),
	}

	return client, nil
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
