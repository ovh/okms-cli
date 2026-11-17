package kmip

import (
	"fmt"

	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/ovh/okms-cli/common/output"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/spf13/cobra"
)

func archiveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "archive ID",
		Short: "Archive an object",
		Long: `Archive an object.

An archived object stays managed by the server but cannot be used until it is recovered.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			resp := exit.OnErr2(kmipClient.Archive(args[0]).ExecContext(cmd.Context()))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
				return
			}
			fmt.Println("Archived object", resp.UniqueIdentifier)
		},
	}
}

func recoverCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "recover ID",
		Short: "Recover an archived object",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			resp := exit.OnErr2(kmipClient.Recover(args[0]).ExecContext(cmd.Context()))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
				return
			}
			fmt.Println("Recovered object", resp.UniqueIdentifier)
		},
	}
}
