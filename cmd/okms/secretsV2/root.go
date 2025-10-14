package secretsv2

import (
	"github.com/ovh/okms-cli/cmd/okms/common"
	"github.com/spf13/cobra"
)

func CreateCommand(cust common.CustomizeFunc) *cobra.Command {
	var kvCmd = &cobra.Command{
		Use:     "secrets",
		Aliases: []string{"kv2", "secret"}, // TODO discuss keywords to use for better UX
		Short:   "This command has subcommands for interacting with KMS's key-value store.",
	}

	common.SetupRestApiFlags(kvCmd, cust)

	kvCmd.AddCommand(
		secretConfigCommand(),
		secretListCmd(),
		secretPostCmd(),
		secretGetCmd(),
		secretPutCmd(),
		secretDeleteCmd(),
		secretVersionCommand(),
	)

	return kvCmd
}
