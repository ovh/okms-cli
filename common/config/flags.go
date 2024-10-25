package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/spf13/cobra"
)

func SetupEndpointFlags(command *cobra.Command, service string, init func(cmd *cobra.Command, cfg EndpointConfig)) {
	service = strings.ToLower(service)
	desc := "KMS endpoint URL"
	if service != "" && service != "http" {
		desc = fmt.Sprintf("Endpoint address to %s", service)
	}
	command.PersistentFlags().String("endpoint", "", desc)
	command.PersistentFlags().String("cert", "", "Path to certificate")
	command.PersistentFlags().String("ca", "", "Path to CA bundle")
	command.PersistentFlags().String("key", "", "Path to key file")
	command.PersistentFlags().Var(new(AuthMethodFlag), "auth-method", "Authentication method to use")
	var format = flagsmgmt.TEXT_OUTPUT_FORMAT
	command.PersistentFlags().Var(&format, "output", "The formatting style for command output.")

	command.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		configFile, _ := cmd.Flags().GetString("config")

		cfg := LoadEndpointConfig(cmd, service, configFile)
		init(cmd, cfg)
	}
}

func SetupConfigFlags(command *cobra.Command) {
	prof := os.Getenv("KMS_PROFILE")
	if prof == "" {
		prof = "default"
	}

	command.PersistentFlags().String("profile", prof, "Name of the profile")
	command.PersistentFlags().StringP("config", "c", "", "Path to a non default configuration file")
	// command.PersistentFlags().BoolP("debug", "d", false, "Activate debug mode")
}
