package common

import (
	"os"

	"github.com/ovh/okms-cli/common/config"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/ovh/okms-sdk-go"
	"github.com/spf13/cobra"
)

var ccmRestClient *okms.Client

func Client() *okms.Client {
	return ccmRestClient
}

type CustomizeFunc func(c *cobra.Command) func(*okms.Client)

func SetupRestApiFlags(command *cobra.Command, cust CustomizeFunc) {
	debug := command.PersistentFlags().BoolP("debug", "d", false, "Activate debug mode")
	retry := command.PersistentFlags().Uint32("retry", 4, "Maximum number of HTTP retries")
	timeout := command.PersistentFlags().Duration("timeout", okms.DefaultHTTPClientTimeout, "Timeout duration for HTTP requests")

	f := func(*okms.Client) {}
	if cust != nil {
		f = cust(command)
	}

	config.SetupEndpointFlags(command, "http", func(command *cobra.Command, cfg config.EndpointConfig) {
		clientCfg := okms.ClientConfig{
			Timeout: timeout,
			TlsCfg:  cfg.TlsConfig(""),
			Retry: &okms.RetryConfig{
				RetryMax: int(*retry),
			},
		}
		if *debug {
			clientCfg.Middleware = okms.DebugTransport(os.Stderr)
		}
		ccmRestClient = exit.OnErr2(okms.NewRestAPIClient(cfg.Endpoint, clientCfg))
		f(ccmRestClient)
	})
}
