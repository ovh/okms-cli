package common

import (
	"os"

	"github.com/google/uuid"
	"github.com/ovh/okms-cli/common/config"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/ovh/okms-sdk-go"
	"github.com/spf13/cobra"
)

var (
	okmsId     uuid.UUID
	restClient *okms.Client
)

func Client() *okms.Client {
	return restClient
}

func GetOkmsId() uuid.UUID {
	return okmsId
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
		okmsId = cfg.Auth.GetOkmsId()
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
		restClient = exit.OnErr2(okms.NewRestAPIClient(cfg.Endpoint, clientCfg))

		if cfg.Auth.GetToken() != nil {
			restClient.SetCustomHeader("Authorization", "Bearer "+*cfg.Auth.GetToken())
		}

		f(restClient)
	})
}
