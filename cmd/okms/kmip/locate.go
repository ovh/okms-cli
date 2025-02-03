package kmip

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/ovh/kmip-go"
	"github.com/ovh/kmip-go/payloads"
	"github.com/ovh/kmip-go/ttlv"
	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/ovh/okms-cli/common/flagsmgmt/kmipflags"
	"github.com/ovh/okms-cli/common/output"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/spf13/cobra"
)

func locateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "locate",
		Aliases: []string{"list", "ls"},
		Short:   "List kmip objects",
		Args:    cobra.NoArgs,
	}

	detailed := cmd.Flags().Bool("details", false, "Display detailed information")
	var state kmipflags.State
	cmd.Flags().Var(&state, "state", "List only object with the given state")
	var objectType kmipflags.ObjectType
	cmd.Flags().Var(&objectType, "type", "List only objects of the given type")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		req := kmipClient.Locate()
		if state != 0 {
			req = req.WithAttribute(kmip.AttributeNameState, kmip.State(state))
		}
		if objectType != 0 {
			req = req.WithObjectType(kmip.ObjectType(objectType))
		}
		locateResp := exit.OnErr2(req.ExecContext(cmd.Context()))
		if !*detailed {
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(locateResp)
			} else {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Id"})
				for _, id := range locateResp.UniqueIdentifier {
					table.Append([]string{id})
				}
				table.Render()
			}
			return
		}

		attributes := []*payloads.GetAttributesResponsePayload{}
		for _, id := range locateResp.UniqueIdentifier {
			attributes = append(attributes, exit.OnErr2(kmipClient.GetAttributes(id).ExecContext(cmd.Context())))
		}
		if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
			output.JsonPrint(attributes)
		} else {
			printObjectTable(attributes)
		}
	}

	return cmd
}

func printObjectTable(objects []*payloads.GetAttributesResponsePayload) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "TYPE", "NAME", "STATE", "ALGORITHM", "SIZE"})
	for _, attr := range objects {
		var row [6]string
		row[0] = attr.UniqueIdentifier
		for _, v := range attr.Attribute {
			if idx := v.AttributeIndex; idx != nil && *idx > 0 {
				continue
			}
			switch v.AttributeName {
			case kmip.AttributeNameObjectType:
				row[1] = ttlv.EnumStr(v.AttributeValue.(kmip.ObjectType))
			case kmip.AttributeNameName:
				row[2] = v.AttributeValue.(kmip.Name).NameValue
			case kmip.AttributeNameState:
				row[3] = ttlv.EnumStr(v.AttributeValue.(kmip.State))
			case kmip.AttributeNameCryptographicAlgorithm:
				row[4] = ttlv.EnumStr(v.AttributeValue.(kmip.CryptographicAlgorithm))
			case kmip.AttributeNameCryptographicLength:
				row[5] = strconv.Itoa(int(v.AttributeValue.(int32)))
			}
		}
		table.Append(row[:])
	}
	table.Render()
}
