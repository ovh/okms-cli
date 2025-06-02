package keys

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/olekukonko/tablewriter"
	"github.com/ovh/okms-cli/cmd/okms/common"
	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/ovh/okms-cli/common/output"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/spf13/cobra"
)

func newDataKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "datakeys",
		Aliases: []string{"datakey", "dk"},
		Short:   "Manage data keys",
	}
	cmd.AddCommand(
		newGenerateDataKeyFromServiceKeyCmd(),
		newDecryptDataKeyCmd(),
	)
	return cmd
}

func newGenerateDataKeyFromServiceKeyCmd() *cobra.Command {
	var (
		keySize int32
		name    string
	)

	cmd := &cobra.Command{
		Use:   "new KEY-ID",
		Short: "Generate data key wrapped by domain key",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			keyId := exit.OnErr2(uuid.Parse(args[0]))
			plaintext, encrypted := exit.OnErr3(common.Client().GenerateDataKey(cmd.Context(), keyId, name, keySize))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(map[string]any{
					"plain":     plaintext,
					"encrypted": encrypted,
				})
			} else {
				fmt.Printf("Plain key: %s\n", base64.StdEncoding.EncodeToString(plaintext))
				fmt.Printf("Encrypted key: %s\n", encrypted)
			}
		},
	}
	cmd.Flags().Int32Var(&keySize, "size", 256, "Size of the key int bits to be generated")
	cmd.Flags().StringVar(&name, "name", "", "Optional name for the data-key")
	return cmd
}

func newDecryptDataKeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "decrypt KEY-ID DATA-KEY",
		Short: "Decrypt data key encrypted by domain key",
		Long: `Decrypt data key encrypted by domain key.

DATA-KEY can be either plain text, a '-' to read from stdin, or a filename prefixed with @`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			keyId := exit.OnErr2(uuid.Parse(args[0]))
			plaintext := exit.OnErr2(common.Client().DecryptDataKey(cmd.Context(), keyId, flagsmgmt.StringFromArg(args[1], 8192)))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(plaintext)
			} else {
				table := tablewriter.NewWriter(os.Stdout)
				exit.OnErr(table.Append([]string{"Plaintext Key", base64.StdEncoding.EncodeToString(plaintext)}))
				exit.OnErr(table.Render())
			}
		},
	}
}
