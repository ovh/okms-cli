package kmip

import (
	"errors"
	"fmt"

	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/ovh/okms-cli/common/output"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/spf13/cobra"
)

func rekeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rekey ID",
		Short: "Rekey a symmetric key object",
		Args:  cobra.ExactArgs(1),
	}

	offset := cmd.Flags().Duration("offset", 0, "Optional rekeying offset")

	// cmd.Flags().Var(&usage, "usage", "Cryptographic usage")
	// name := cmd.Flags().String("name", "", "Optional key name")

	// sensitive := cmd.Flags().Bool("sensitive", false, "Change sensitive attribute")
	// extractable := cmd.Flags().Bool("extractable", false, "Change the extractable attribute")
	// description := cmd.Flags().String("description", "", "Change the description attribute")
	// comment := cmd.Flags().String("comment", "", "Change the comment attribute")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		req := kmipClient.Rekey(args[0])

		if cmd.Flag("offset").Changed {
			if *offset < 0 {
				exit.OnErr(errors.New("offset cannot be negative"))
			}
			req = req.WithOffset(*offset)
		}
		// if cmd.Flag("sensitive").Changed {
		// 	req = req.WithAttribute(kmip.AttributeNameSensitive, *sensitive)
		// }
		// if cmd.Flag("extractable").Changed {
		// 	req = req.WithAttribute(kmip.AttributeNameExtractable, *extractable)
		// }
		// if *description != "" {
		// 	req = req.WithAttribute(kmip.AttributeNameDescription, *description)
		// }
		// if *comment != "" {
		// 	req = req.WithAttribute(kmip.AttributeNameComment, *comment)
		// }

		resp := exit.OnErr2(req.ExecContext(cmd.Context()))
		if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
			output.JsonPrint(resp)
		} else {
			fmt.Println("Replacement key ID:", resp.UniqueIdentifier)
			if attr := resp.TemplateAttribute; attr != nil && len(attr.Attribute) > 0 {
				printAttributeTable(attr.Attribute)
			}
		}
	}

	return cmd
}
