package keys

import (
	"github.com/ovh/okms-cli/cmd/okms/common"
	"github.com/spf13/cobra"
)

func CreateCommand(cust common.CustomizeFunc) *cobra.Command {
	keysCmd := &cobra.Command{
		Use:     "keys",
		Aliases: []string{"key"},
		Short:   "Manage domain keys",
	}

	common.SetupRestApiFlags(keysCmd, cust)

	keysCmd.AddCommand(
		newListServiceKeysCmd(),
		newAddServiceKeyCmd(),
		newUpdateKeyCmd(),
		newExportPublicKeyCmd(),
		newImportServiceKeyCmd(),
		newGetServiceKeyCmd(),
		newEncryptWithServiceKeyCmd(),
		newDecryptWithServiceKeyCmd(),
		newDataKeysCmd(),
		newSignCmd(),
		newVerifyCmd(),
		newDeactivateKeyCmd(),
		newDeleteKeyCmd(),
		newActivateKeyCmd(),
	)

	return keysCmd
}
