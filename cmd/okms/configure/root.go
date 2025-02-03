package configure

import (
	"fmt"

	"github.com/ovh/okms-cli/common/config"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func CreateCommand() *cobra.Command {
	configureCmd := &cobra.Command{
		Use:     "configure",
		Aliases: []string{"config"},
		Short:   "Configure CLI options",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			profile := cmd.Flag("profile").Value.String()
			configFile := cmd.Flag("config").Value.String()

			file, _ := config.LoadFromFile("okms", configFile)

			if mig, _ := cmd.Flags().GetBool("migrate-only"); mig {
				exit.OnErr(config.WriteToFile(file))
				fmt.Println("Migration completed")
				return
			}

			Run(profile)
			if exit.OnErr2(pterm.DefaultInteractiveConfirm.Show("Save the new configuration ?")) {
				exit.OnErr(config.WriteToFile(file))
			}
		},
	}

	configureCmd.Flags().Bool("migrate-only", false, "Migrate to latest config schema, without changing the settings")

	configureCmd.AddCommand(config.SetProfileCommand("okms"))

	return configureCmd
}

func Run(profile string) {
	choice := exit.OnErr2(pterm.DefaultInteractiveSelect.WithOptions([]string{"HTTP", "KMIP"}).Show("Select a protocol to configure"))
	if choice == "HTTP" {
		config.ReadUserInput("CA file", "http.ca", profile, config.ValidateFileExists.AllowEmpty())
		config.ReadUserInput("Certificate file", "http.auth.cert", profile, config.ValidateFileExists)
		config.ReadUserInput("Private key file", "http.auth.key", profile, config.ValidateFileExists)
		config.ReadUserInput("Endpoint", "http.endpoint", profile, config.ValidateURL)
	} else if choice == "KMIP" {
		config.ReadUserInput("CA file", "kmip.ca", profile, config.ValidateFileExists.AllowEmpty())
		config.ReadUserInput("Certificate file", "kmip.auth.cert", profile, config.ValidateFileExists)
		config.ReadUserInput("Private key file", "kmip.auth.key", profile, config.ValidateFileExists)
		config.ReadUserInput("Endpoint", "kmip.endpoint", profile, config.ValidateTCPAddr)
	}
}
