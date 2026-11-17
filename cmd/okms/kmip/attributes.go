package kmip

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/ovh/kmip-go"
	"github.com/ovh/kmip-go/ttlv"
	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/ovh/okms-cli/common/output"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/spf13/cobra"
)

var (
	attributeValueHdrRegex    = regexp.MustCompile(`^AttributeValue \(.+\): `)
	attributeValueFieldsRegex = regexp.MustCompile(`(.+) \(.+\): `)
)

func printAttributeTable(attributes []kmip.Attribute) {
	opts := tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
		Borders: tw.Border{Left: tw.On, Right: tw.On, Top: tw.On, Bottom: tw.On},
		Settings: tw.Settings{
			Separators: tw.Separators{ShowHeader: tw.On, BetweenRows: tw.On, BetweenColumns: tw.On},
		},
	}))
	table := tablewriter.NewTable(os.Stdout, opts).Configure(func(c *tablewriter.Config) {
		c.Header.Formatting.AutoWrap = tw.WrapNormal
		c.Header.Alignment.Global = tw.AlignLeft
		c.Row.Formatting.AutoWrap = tw.WrapNormal
		c.Row.Alignment.Global = tw.AlignLeft
	})

	enc := ttlv.NewTextEncoder()
	for _, attr := range attributes {
		enc.Clear()
		enc.TagAny(kmip.TagAttributeValue, attr.AttributeValue)
		txt := enc.Bytes()

		txt = attributeValueHdrRegex.ReplaceAll(txt, nil)
		txt = attributeValueFieldsRegex.ReplaceAll(txt, []byte("$1: "))
		txt = bytes.ReplaceAll(txt, []byte("\n    "), []byte("\n"))
		txt = bytes.TrimSpace(txt)

		name := string(attr.AttributeName)
		if idx := attr.AttributeIndex; idx != nil && *idx > 0 {
			name = fmt.Sprintf("%s [%d]", name, *idx)
		}
		exit.OnErr(table.Append([]string{name, string(txt)}))
	}

	exit.OnErr(table.Render())
}

func getAttributesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get ID",
		Short: "Get the attributes of an object",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			attributes := exit.OnErr2(kmipClient.GetAttributes(args[0]).ExecContext(cmd.Context()))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(attributes)
				return
			}
			printAttributeTable(attributes.Attribute)
		},
	}
}

func deleteAttributeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete ID ATTRIBUTE_NAME",
		Short: "Delete an existing attribute of an object",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]
			attrName := kmip.AttributeName(args[1])

			req := kmipClient.DeleteAttribute(id, attrName)
			if idx, err := cmd.Flags().GetInt32("index"); err == nil && cmd.Flags().Changed("index") {
				req = req.WithIndex(idx)
			}

			resp := exit.OnErr2(req.ExecContext(cmd.Context()))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
				return
			}
			printAttributeTable([]kmip.Attribute{resp.Attribute})
		},
	}
	cmd.Flags().Int32("index", 0, "Index of the attribute instance to delete (default 0)")
	return cmd
}

func modifyAttributeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify ID ATTRIBUTE_NAME VALUE",
		Short: "Modify an existing attribute of an object",
		Long: `Modify an existing attribute of a KMIP managed object.

For the "Name" attribute, VALUE is passed as an 'Uninterpreted Text String'.
For all other standard attributes, VALUE is passed as a plain 'Text String'.
Therefore, only attributes with a 'Text String' encoding are supported with this command.`,
		Args: cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]
			attrName := kmip.AttributeName(args[1])
			value := textAttributeValue(attrName, args[2])

			resp := exit.OnErr2(kmipClient.ModifyAttribute(id, attrName, value).ExecContext(cmd.Context()))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
				return
			}
			printAttributeTable([]kmip.Attribute{resp.Attribute})
		},
	}
	return cmd
}

func addAttributeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add ID ATTRIBUTE_NAME VALUE",
		Short: "Add a new attribute to an object",
		Long: `Add a new attribute to a KMIP managed object.

For the "Name" attribute, VALUE is passed as an 'Uninterpreted Text String'.
For all other standard attributes, VALUE is passed as a plain 'Text String'.
Therefore, only attributes with a 'Text String' encoding are supported with this command.`,
		Args: cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]
			attrName := kmip.AttributeName(args[1])
			value := textAttributeValue(attrName, args[2])

			req := kmipClient.AddAttribute(id, attrName, value)
			if idx, err := cmd.Flags().GetInt32("index"); err == nil && cmd.Flags().Changed("index") {
				req = req.WithIndex(idx)
			}

			resp := exit.OnErr2(req.ExecContext(cmd.Context()))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
				return
			}
			printAttributeTable([]kmip.Attribute{resp.Attribute})
		},
	}
	cmd.Flags().Int32("index", 0, "Index of the attribute instance to add (default 0)")
	return cmd
}

func listAttributesCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list ID",
		Aliases: []string{"ls"},
		Short:   "List the names of the attributes of an object",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			resp := exit.OnErr2(kmipClient.GetAttributeList(args[0]).ExecContext(cmd.Context()))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
				return
			}
			for _, name := range resp.AttributeName {
				fmt.Println(name)
			}
		},
	}
}

// textAttributeValue converts the VALUE argument of the attributes add and modify commands
// into an attribute value: a Name for the "Name" attribute, a text string otherwise.
func textAttributeValue(name kmip.AttributeName, rawValue string) any {
	if name == kmip.AttributeNameName {
		return kmip.Name{
			NameValue: rawValue,
			NameType:  kmip.NameTypeUninterpretedTextString,
		}
	}
	return rawValue
}

func attributesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "attributes",
		Aliases: []string{"attribute", "attr"},
		Short:   "Manage an object's attributes",
	}
	cmd.AddCommand(
		getAttributesCommand(),
		listAttributesCommand(),
		addAttributeCommand(),
		modifyAttributeCommand(),
		deleteAttributeCommand(),
	)
	return cmd
}
