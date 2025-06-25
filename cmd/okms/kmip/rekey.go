package kmip

import (
	"errors"
	"fmt"
	"time"

	"github.com/ovh/kmip-go"
	"github.com/ovh/kmip-go/ttlv"
	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/ovh/okms-cli/common/output"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/spf13/cobra"
)

func rekeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rekey ID",
		Short: "Rekey a symmetric key or a key-pair",
		Args:  cobra.ExactArgs(1),
	}

	offset := cmd.Flags().Duration("offset", 0, "Optional rekeying offset")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		attrs := exit.OnErr2(kmipClient.GetAttributes(args[0], kmip.AttributeNameObjectType).ExecContext(cmd.Context()))
		var objType kmip.ObjectType
		for _, attr := range attrs.Attribute {
			if attr.AttributeName == kmip.AttributeNameObjectType {
				objType = attr.AttributeValue.(kmip.ObjectType)
			}
		}
		if objType == 0 {
			exit.Now("Missing object type from server returned attributes")
		}

		switch objType {
		case kmip.ObjectTypeSymmetricKey:
			rekeySymmetric(cmd, args, offset)
		case kmip.ObjectTypePrivateKey:
			rekeyKeypair(cmd, args, offset)
		case kmip.ObjectTypePublicKey:
			exit.Now("Cannot rekey public-key, please specify the private-key ID instead")
		default:
			exit.Now("Cannot rekey an object of type %s", ttlv.EnumStr(objType))
		}
	}

	return cmd
}

func rekeySymmetric(cmd *cobra.Command, args []string, offset *time.Duration) {
	req := kmipClient.Rekey(args[0])
	if cmd.Flag("offset").Changed {
		if *offset < 0 {
			exit.OnErr(errors.New("offset cannot be negative"))
		}
		req = req.WithOffset(*offset)
	}

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

func rekeyKeypair(cmd *cobra.Command, args []string, offset *time.Duration) {
	req := kmipClient.RekeyKeyPair(args[0])
	if cmd.Flag("offset").Changed {
		if *offset < 0 {
			exit.OnErr(errors.New("offset cannot be negative"))
		}
		req = req.WithOffset(*offset)
	}

	resp := exit.OnErr2(req.ExecContext(cmd.Context()))
	if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
		output.JsonPrint(resp)
	} else {
		fmt.Println("Replacement private-key ID:", resp.PrivateKeyUniqueIdentifier)
		fmt.Println("Replacement public-key ID:", resp.PublicKeyUniqueIdentifier)
	}
}
