package secrets

import (
	"github.com/ovh/okms-cli/cmd/okms/common"
	"github.com/spf13/cobra"
)

func CreateCommand(cust common.CustomizeFunc) *cobra.Command {
	var kvCmd = &cobra.Command{
		Use:     "secrets",
		Aliases: []string{"kv", "secret"},
		Short:   "This command has subcommands for interacting with KMS's key-value store.",
	}

	common.SetupRestApiFlags(kvCmd, cust)

	kvCmd.AddCommand(
		kvGetCmd(),
		kvPutCmd(),
		kvPatchCmd(),
		kvDeleteCmd(),
		kvUndeleteCmd(),
		kvDestroyCmd(),
		kvConfigCommand(),
		kvSubkeysCmd(),
		kvMetadataCommand(),
	)

	return kvCmd
}
