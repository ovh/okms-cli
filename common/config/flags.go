package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/spf13/cobra"
)

// endpointInitializer holds what is needed to initialize the endpoint of a
// command tree: the service name used to read the configuration, and the
// callback building the client.
type endpointInitializer struct {
	service string
	init    func(*cobra.Command, EndpointConfig)
}

// endpointInitializers holds the initializer registered by SetupEndpointFlags
// for each command owning an endpoint. It is used to initialize the endpoint
// outside of the command lifecycle, typically from shell completion functions.
var endpointInitializers = map[*cobra.Command]endpointInitializer{}

func SetupEndpointFlags(command *cobra.Command, service string, init func(cmd *cobra.Command, cfg EndpointConfig)) {
	service = strings.ToLower(service)
	desc := "KMS endpoint URL"
	if service != "" && service != "http" && service != "restapi" {
		desc = fmt.Sprintf("Endpoint address to %s", service)
	}
	command.PersistentFlags().String("endpoint", "", desc)
	command.PersistentFlags().String("ca", "", "Path to CA bundle")
	command.PersistentFlags().Bool("insecure-skip-verify", false, "Skip server TLS certificate verification (insecure)")
	command.PersistentFlags().String("cert", "", "Path to certificate")
	command.PersistentFlags().String("key", "", "Path to key file")
	command.PersistentFlags().String("token", "", "Token")
	command.PersistentFlags().String("okmsId", "", "OKMS id")
	command.PersistentFlags().Var(new(AuthMethodFlag), "auth-method", "Authentication method to use")

	var format = flagsmgmt.TEXT_OUTPUT_FORMAT
	command.PersistentFlags().Var(&format, "output", "The formatting style for command output.")

	initializer := endpointInitializer{service: service, init: init}
	endpointInitializers[command] = initializer

	command.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		initEndpoint(cmd, initializer)
	}
}

// completionMode is set while an endpoint is being initialized from a shell
// completion function, to prevent side effects such as creating a configuration
// file.
var completionMode bool

// InitEndpoint initializes the endpoint of the command tree cmd belongs to
// (for instance by building the API client), without going through the command
// lifecycle. Shell completion functions run before the PersistentPreRun hooks,
// so they must call InitEndpoint to get a usable client.
func InitEndpoint(cmd *cobra.Command) error {
	initializer, ok := findEndpointInitializer(cmd)
	if !ok {
		return fmt.Errorf("no endpoint initializer found for command %q", cmd.CommandPath())
	}
	completionMode = true
	defer func() { completionMode = false }()
	initEndpoint(cmd, initializer)
	return nil
}

// findEndpointInitializer returns the initializer of the closest command
// owning an endpoint, starting from cmd and walking up its parents.
func findEndpointInitializer(cmd *cobra.Command) (endpointInitializer, bool) {
	for c := cmd; c != nil; c = c.Parent() {
		if initializer, ok := endpointInitializers[c]; ok {
			return initializer, true
		}
	}
	return endpointInitializer{}, false
}

func initEndpoint(cmd *cobra.Command, initializer endpointInitializer) {
	configFile, _ := cmd.Flags().GetString("config")

	cfg := LoadEndpointConfig(cmd, initializer.service, configFile)
	initializer.init(cmd, cfg)
}

func SetupConfigFlags(command *cobra.Command) {
	prof := os.Getenv("KMS_PROFILE")
	if prof == "" {
		prof = "default"
	}

	command.PersistentFlags().String("profile", prof, "Name of the profile")
	command.PersistentFlags().StringP("config", "c", "", "Path to a non default configuration file")
	// command.PersistentFlags().BoolP("debug", "d", false, "Activate debug mode")

	_ = command.RegisterFlagCompletionFunc("profile", CompleteProfiles)
}
