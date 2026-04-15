package secrets

import (
	"github.com/ovh/okms-cli/cmd/okms/common"
	"github.com/spf13/cobra"
)

func CreateCommand(cust common.CustomizeFunc) *cobra.Command {
	var kvCmd = &cobra.Command{
		Use:     "vault",
		Aliases: []string{"kv"},
		Short:   "Manage secrets through Hashicorp Vault API",
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
