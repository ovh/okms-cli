package secrets

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/ovh/okms-cli/cmd/okms/common"
	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/ovh/okms-cli/common/output"
	"github.com/ovh/okms-cli/common/utils"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/ovh/okms-sdk-go/types"
	"github.com/spf13/cobra"
)

func kvConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manages secret engine configuration",
	}

	cmd.AddCommand(
		kvReadConfigCommand(),
		kvWriteConfigCommand(),
	)
	return cmd
}

func kvReadConfigCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "read",
		Short: "Reads secret engine configuration",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			resp := exit.OnErr2(common.Client().GetSecretConfig(cmd.Context()))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else {
				table := tablewriter.NewWriter(os.Stdout)
				table.AppendBulk([][]string{
					{"cas", fmt.Sprintf("%t", utils.DerefOrDefault(resp.Data.CasRequired))},
					{"Delete version after", utils.DerefOrDefault(resp.Data.DeleteVersionAfter)},
					{"Max. number of versions", fmt.Sprintf("%d", utils.DerefOrDefault(resp.Data.MaxVersions))},
				})
				table.Render()
			}
		},
	}
}

func kvWriteConfigCommand() *cobra.Command {
	var (
		casRequired        bool
		maxVersions        uint32
		deleteVersionAfter string
	)

	cmd := &cobra.Command{
		Use:   "write",
		Short: "Writes secret engine configuration",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var c *bool
			if cmd.Flag("cas-required").Changed {
				c = &casRequired
			}

			var d *string
			if cmd.Flag("delete-after").Changed {
				d = &deleteVersionAfter
			}

			var m *uint32
			if cmd.Flag("max-versions").Changed {
				m = &maxVersions
			}

			body := types.PostConfigRequest{
				CasRequired:        c,
				DeleteVersionAfter: d,
				MaxVersions:        m,
			}

			exit.OnErr(common.Client().PostSecretConfig(cmd.Context(), body))
		},
	}

	cmd.Flags().BoolVar(&casRequired, "cas-required", false, "If true all keys will require the cas parameter to be set on all write requests.")
	cmd.Flags().Uint32Var(&maxVersions, "max-versions", 0, "The number of versions to keep per key. This value applies to all keys, but a key's metadata setting can overwrite this value. Once a key has more than the configured allowed versions, the oldest version will be permanently deleted. ")
	cmd.Flags().StringVar(&deleteVersionAfter, "delete-after", "0s", "If set, specifies the length of time before a version is deleted.\nDate format, see: https://developer.hashicorp.com/vault/docs/concepts/duration-format")
	return cmd
}
