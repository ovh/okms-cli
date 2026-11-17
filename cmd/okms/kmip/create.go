package kmip

import (
	"errors"
	"fmt"

	"github.com/ovh/kmip-go"
	"github.com/ovh/kmip-go/kmipclient"
	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/ovh/okms-cli/common/flagsmgmt/kmipflags"
	"github.com/ovh/okms-cli/common/output"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/spf13/cobra"
)

func createCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create kmip keys",
	}
	cmd.AddCommand(
		createSymmetricKey(),
		createKeyPair(),
	)
	return cmd
}

func createSymmetricKey() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "symmetric",
		Aliases: []string{"sym"},
		Short:   "Create KMIP symmetric key",
		Args:    cobra.NoArgs,
	}

	var alg kmipflags.SymmetricAlg
	usage := kmipflags.KeyUsageList{kmipflags.ENCRYPT, kmipflags.DECRYPT}

	cmd.Flags().Var(&alg, "alg", "Key algorithm")
	size := cmd.Flags().Int("size", 0, "Key bit length")
	cmd.Flags().Var(&usage, "usage", "Cryptographic usage")
	name := cmd.Flags().String("name", "", "Optional key name")

	sensitive := cmd.Flags().Bool("sensitive", false, "Set sensitive attribute")
	extractable := cmd.Flags().Bool("extractable", true, "Set the extractable attribute")
	description := cmd.Flags().String("description", "", "Set the description attribute")
	comment := cmd.Flags().String("comment", "", "Set the comment attribute")
	usageLimits := setUsageLimitsFlags(cmd)

	_ = cmd.MarkFlagRequired("alg")
	_ = cmd.MarkFlagRequired("size")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		req := kmipClient.Create().
			SymmetricKey(kmip.CryptographicAlgorithm(alg), *size, usage.ToCryptographicUsageMask())
		if cmd.Flags().Changed("extractable") {
			req = req.WithAttribute(kmip.AttributeNameExtractable, *extractable)
		}
		if cmd.Flags().Changed("sensitive") {
			req = req.WithAttribute(kmip.AttributeNameSensitive, *sensitive)
		}
		if *name != "" {
			req = req.WithName(*name)
		}
		if *description != "" {
			req = req.WithAttribute(kmip.AttributeNameDescription, *description)
		}
		if *comment != "" {
			req = req.WithAttribute(kmip.AttributeNameComment, *comment)
		}
		req = withUsageLimits(cmd, usageLimits, req)

		resp := exit.OnErr2(req.ExecContext(cmd.Context()))

		if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
			output.JsonPrint(resp)
		} else {
			fmt.Println("Key created with ID", resp.UniqueIdentifier)
			// Print returned attributes if any
			if resp.Attributes != nil && len(resp.Attributes.Attribute) > 0 {
				printAttributeTable(resp.Attributes.Attribute)
			}
		}
	}

	return cmd
}

func createKeyPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key-pair",
		Short: "Create an asymmetric key-pair",
		Args:  cobra.NoArgs,
	}

	var alg kmipflags.AsymmetricAlg
	cmd.Flags().Var(&alg, "alg", "Key-pair algorithm")
	size := cmd.Flags().Int("size", 0, "Modulus  bit length of the RSA key-pair to generate")
	var curve kmipflags.EcCurve
	cmd.Flags().Var(&curve, "curve", "Elliptic curve for EC keys")
	privateUsage := kmipflags.KeyUsageList{kmipflags.SIGN}
	publicUsage := kmipflags.KeyUsageList{kmipflags.VERIFY}
	cmd.Flags().Var(&privateUsage, "private-usage", "Private key allowed usage")
	cmd.Flags().Var(&publicUsage, "public-usage", "Public key allowed usage")
	privateName := cmd.Flags().String("private-name", "", "Optional private key name")
	publicName := cmd.Flags().String("public-name", "", "Optional public key name")

	privateSensitive := cmd.Flags().Bool("private-sensitive", false, "Set sensitive attribute on the private key")
	privateExtractable := cmd.Flags().Bool("private-extractable", true, "Set the extractable attribute on the private key")
	description := cmd.Flags().String("description", "", "Set the description attribute on both keys")
	comment := cmd.Flags().String("comment", "", "Set the comment attribute on both keys")
	usageLimits := setUsageLimitsFlags(cmd)

	_ = cmd.MarkFlagRequired("alg")
	cmd.MarkFlagsMutuallyExclusive("curve", "size")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		var req kmipclient.ExecCreateKeyPairAttr
		switch alg {
		case kmipflags.RSA:
			if *size == 0 {
				exit.OnErr(errors.New("Missing --size flag"))
			}
			req = kmipClient.CreateKeyPair().RSA(*size, privateUsage.ToCryptographicUsageMask(), publicUsage.ToCryptographicUsageMask())
		case kmipflags.ECDSA:
			if curve == 0 {
				exit.OnErr(errors.New("Missing --curve flag"))
			}
			req = kmipClient.CreateKeyPair().ECDSA(kmip.RecommendedCurve(curve), privateUsage.ToCryptographicUsageMask(), publicUsage.ToCryptographicUsageMask())
		}
		if *privateName != "" {
			req = req.PrivateKey().WithName(*privateName)
		}
		if *publicName != "" {
			req = req.PublicKey().WithName(*publicName)
		}
		req = req.PrivateKey()
		if cmd.Flags().Changed("private-extractable") {
			req = req.WithAttribute(kmip.AttributeNameExtractable, *privateExtractable)
		}
		if cmd.Flags().Changed("private-sensitive") {
			req = req.WithAttribute(kmip.AttributeNameSensitive, *privateSensitive)
		}
		if *description != "" {
			req = req.Common().WithAttribute(kmip.AttributeNameDescription, *description)
		}
		if *comment != "" {
			req = req.Common().WithAttribute(kmip.AttributeNameComment, *comment)
		}
		req = withUsageLimits(cmd, usageLimits, req.Common())
		resp := exit.OnErr2(req.ExecContext(cmd.Context()))

		if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
			output.JsonPrint(resp)
		} else {
			fmt.Println("Pubic Key ID:", resp.PublicKeyUniqueIdentifier)
			fmt.Println("Private Key ID:", resp.PrivateKeyUniqueIdentifier)
			// Print returned attributes if any
			if attrs := resp.PublicKeyTemplateAttribute; attrs != nil && len(attrs.Attribute) > 0 {
				fmt.Println("Public Key Attributes:")
				printAttributeTable(attrs.Attribute)
			}
			if attrs := resp.PrivateKeyTemplateAttribute; attrs != nil && len(attrs.Attribute) > 0 {
				fmt.Println("Private Key Attributes:")
				printAttributeTable(attrs.Attribute)
			}
		}
	}

	return cmd
}

type usageLimitsFlags struct {
	total int64
	unit  kmipflags.UsageLimitsUnit
}

func setUsageLimitsFlags(cmd *cobra.Command) *usageLimitsFlags {
	f := &usageLimitsFlags{unit: kmipflags.UsageLimitsObject}
	cmd.Flags().Int64Var(&f.total, "usage-limit", 0, "Set the usage limits attribute: total amount of usage allowed for the key, in --usage-limit-unit")
	cmd.Flags().Var(&f.unit, "usage-limit-unit", "Unit of --usage-limit: number of cryptographic operations (Object) or of processed bytes (Byte)")
	return f
}

// withUsageLimits adds the Usage Limits attribute to the request when --usage-limit is set.
// The remaining usage count starts at the total.
func withUsageLimits[T interface {
	WithAttribute(kmip.AttributeName, any) T
}](cmd *cobra.Command, f *usageLimitsFlags, req T) T {
	if !cmd.Flags().Changed("usage-limit") {
		return req
	}
	if f.total <= 0 {
		exit.OnErr(errors.New("--usage-limit must be greater than 0"))
	}
	count := f.total
	return req.WithAttribute(kmip.AttributeNameUsageLimits, kmip.UsageLimits{
		UsageLimitsTotal: f.total,
		UsageLimitsCount: &count,
		UsageLimitsUnit:  kmip.UsageLimitsUnit(f.unit),
	})
}
