package kmip

import (
	"encoding/base64"
	"fmt"

	"github.com/ovh/kmip-go"
	"github.com/ovh/kmip-go/kmipclient"
	"github.com/ovh/kmip-go/payloads"
	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/ovh/okms-cli/common/flagsmgmt/kmipflags"
	"github.com/ovh/okms-cli/common/output"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/spf13/cobra"
)

func registerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a kmip object",
	}

	cmd.AddCommand(
		registerSecretCommand(),
		registerSymmetricKeyCommand(),
		registerCertificateCommand(),
		registerPublicKeyCommand(),
		registerPrivateKeyCommand(),
		registerKeyPairCommand(),
	)

	return cmd
}

func registerSecretCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret VALUE",
		Short: "Register a secret object",
		Long: `Register a secret object.

VALUE can be either plain text, a '-' to read from stdin, or a filename prefixed with @.`,
		Args: cobra.ExactArgs(1),
	}

	b64 := cmd.Flags().Bool("base64", false, "Given secret is base64 encoded")
	name := cmd.Flags().String("name", "", "Optional name for the secret")
	description := cmd.Flags().String("description", "", "Set the description attribute")
	comment := cmd.Flags().String("comment", "", "Set the comment attribute")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		secret := flagsmgmt.BytesFromArg(args[0], 16_000)
		if *b64 {
			secret = exit.OnErr2(base64.StdEncoding.AppendDecode(nil, secret))
		}
		//TODO: Make secret type a flag argument
		req := kmipClient.Register().Secret(kmip.SecretDataTypePassword, secret)
		if *name != "" {
			req = req.WithName(*name)
		}
		if *description != "" {
			req = req.WithAttribute(kmip.AttributeNameDescription, *description)
		}
		if *comment != "" {
			req = req.WithAttribute(kmip.AttributeNameComment, *comment)
		}
		resp := exit.OnErr2(req.ExecContext(cmd.Context()))

		if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
			output.JsonPrint(resp)
		} else {
			fmt.Println("Secret registered with ID", resp.UniqueIdentifier)
			// Print returned attributes if any
			if attr := resp.TemplateAttribute; attr != nil && len(attr.Attribute) > 0 {
				printAttributeTable(attr.Attribute)
			}
		}
	}

	return cmd
}

func registerSymmetricKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "symmetric VALUE",
		Aliases: []string{"sym"},
		Short:   "Register a symmetric key object",
		Long: `Register a symmetric key object.

VALUE can be either plain text, a '-' to read from stdin, or a filename prefixed with @.`,
		Args: cobra.ExactArgs(1),
	}

	var alg kmipflags.SymmetricAlg
	cmd.Flags().Var(&alg, "alg", "Key's cryptographic algorithm")

	usage := kmipflags.KeyUsageList{kmipflags.ENCRYPT, kmipflags.DECRYPT}
	cmd.Flags().Var(&usage, "usage", "Cryptographic usage")

	b64 := cmd.Flags().Bool("base64", false, "Given key is base64 encoded")
	name := cmd.Flags().String("name", "", "Optional name for the key")
	description := cmd.Flags().String("description", "", "Set the description attribute")
	comment := cmd.Flags().String("comment", "", "Set the comment attribute")
	sensitive := cmd.Flags().Bool("sensitive", false, "Set sensitive attribute")
	extractable := cmd.Flags().Bool("extractable", true, "Set the extractable attribute")

	_ = cmd.MarkFlagRequired("alg")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		key := flagsmgmt.BytesFromArg(args[0], 16_000)
		if *b64 {
			key = exit.OnErr2(base64.StdEncoding.AppendDecode(nil, key))
		}

		req := kmipClient.Register().
			SymmetricKey(kmip.CryptographicAlgorithm(alg), usage.ToCryptographicUsageMask(), key).
			WithAttribute(kmip.AttributeNameExtractable, *extractable).
			WithAttribute(kmip.AttributeNameSensitive, *sensitive)
		if *name != "" {
			req = req.WithName(*name)
		}
		if *description != "" {
			req = req.WithAttribute(kmip.AttributeNameDescription, *description)
		}
		if *comment != "" {
			req = req.WithAttribute(kmip.AttributeNameComment, *comment)
		}
		resp := exit.OnErr2(req.ExecContext(cmd.Context()))

		if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
			output.JsonPrint(resp)
		} else {
			fmt.Println("Symmetric key registered with ID", resp.UniqueIdentifier)
			// Print returned attributes if any
			if attr := resp.TemplateAttribute; attr != nil && len(attr.Attribute) > 0 {
				printAttributeTable(attr.Attribute)
			}
		}
	}

	return cmd
}

func registerCertificateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "certificate VALUE",
		Aliases: []string{"cert", "crt"},
		Short:   "Register an X509 certificate",
		Long: `Register an X509 certificate.

VALUE can be either plain text, a '-' to read from stdin, or a filename prefixed with @.`,
		Args: cobra.ExactArgs(1),
	}

	isPem := cmd.Flags().Bool("pem", false, "Certificate is PEM encoded")
	name := cmd.Flags().String("name", "", "Optional name for the certificate")
	description := cmd.Flags().String("description", "", "Set the description attribute")
	comment := cmd.Flags().String("comment", "", "Set the comment attribute")
	publicKeyId := cmd.Flags().String("public-key", "", "Set a link to the certificates public key")
	parent := cmd.Flags().String("parent", "", "Set a link to the parent signing certificate")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		cert := flagsmgmt.BytesFromArg(args[0], 16_000)

		var req kmipclient.ExecRegister
		if *isPem {
			req = kmipClient.Register().PemCertificate(cert)
		} else {
			req = kmipClient.Register().Certificate(kmip.CertificateTypeX_509, cert)
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
		if *publicKeyId != "" {
			req = req.WithLink(kmip.LinkTypePublicKeyLink, *publicKeyId)
		}
		if *parent != "" {
			req = req.WithLink(kmip.LinkTypeCertificateLink, *parent)
		}
		resp := exit.OnErr2(req.ExecContext(cmd.Context()))

		if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
			output.JsonPrint(resp)
		} else {
			fmt.Println("Certificate registered with ID", resp.UniqueIdentifier)
			// Print returned attributes if any
			if attr := resp.TemplateAttribute; attr != nil && len(attr.Attribute) > 0 {
				printAttributeTable(attr.Attribute)
			}
		}
	}

	return cmd
}

func registerPublicKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "public-key VALUE",
		Aliases: []string{"public", "pub"},
		Short:   "Register a public key object from PEM encoded data",
		Long: `Register a public key object from PEM encoded data.

VALUE can be either plain text, a '-' to read from stdin, or a filename prefixed with @.`,
		Args: cobra.ExactArgs(1),
	}
	usage := kmipflags.KeyUsageList{kmipflags.VERIFY}
	cmd.Flags().Var(&usage, "usage", "Cryptographic usage")

	name := cmd.Flags().String("name", "", "Optional name for the key")
	description := cmd.Flags().String("description", "", "Set the description attribute")
	comment := cmd.Flags().String("comment", "", "Set the comment attribute")

	privLink := cmd.Flags().String("private-link", "", "Optional private key ID to link to")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		key := flagsmgmt.BytesFromArg(args[0], 16_000)

		req := kmipClient.Register().PemPublicKey(key, usage.ToCryptographicUsageMask())
		if *name != "" {
			req = req.WithName(*name)
		}
		if *description != "" {
			req = req.WithAttribute(kmip.AttributeNameDescription, *description)
		}
		if *comment != "" {
			req = req.WithAttribute(kmip.AttributeNameComment, *comment)
		}
		if *privLink != "" {
			req = req.WithLink(kmip.LinkTypePrivateKeyLink, *privLink)
		}
		resp := exit.OnErr2(req.ExecContext(cmd.Context()))

		if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
			output.JsonPrint(resp)
		} else {
			fmt.Println("Public key registered with ID", resp.UniqueIdentifier)
			// Print returned attributes if any
			if attr := resp.TemplateAttribute; attr != nil && len(attr.Attribute) > 0 {
				printAttributeTable(attr.Attribute)
			}
		}
	}

	return cmd
}

func registerPrivateKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "private-key VALUE",
		Aliases: []string{"private", "priv"},
		Short:   "Register a private key object from PEM encoded data",
		Long: `Register a private key object from PEM encoded data.
		
VALUE can be either plain text, a '-' to read from stdin, or a filename prefixed with @.`,
		Args: cobra.ExactArgs(1),
	}
	usage := kmipflags.KeyUsageList{kmipflags.SIGN}
	cmd.Flags().Var(&usage, "usage", "Cryptographic usage")

	name := cmd.Flags().String("name", "", "Optional name for the key")
	description := cmd.Flags().String("description", "", "Set the description attribute")
	comment := cmd.Flags().String("comment", "", "Set the comment attribute")

	sensitive := cmd.Flags().Bool("sensitive", false, "Set sensitive attribute")
	extractable := cmd.Flags().Bool("extractable", true, "Set the extractable attribute")

	pubLink := cmd.Flags().String("public-link", "", "Optional public key ID to link to")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		key := flagsmgmt.BytesFromArg(args[0], 16_000)

		req := kmipClient.Register().PemPrivateKey(key, usage.ToCryptographicUsageMask()).
			WithAttribute(kmip.AttributeNameExtractable, *extractable).
			WithAttribute(kmip.AttributeNameSensitive, *sensitive)
		if *name != "" {
			req = req.WithName(*name)
		}
		if *description != "" {
			req = req.WithAttribute(kmip.AttributeNameDescription, *description)
		}
		if *comment != "" {
			req = req.WithAttribute(kmip.AttributeNameComment, *comment)
		}
		if *pubLink != "" {
			req = req.WithLink(kmip.LinkTypePublicKeyLink, *pubLink)
		}
		resp := exit.OnErr2(req.ExecContext(cmd.Context()))

		if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
			output.JsonPrint(resp)
		} else {
			fmt.Println("Private key registered with ID", resp.UniqueIdentifier)
			// Print returned attributes if any
			if attr := resp.TemplateAttribute; attr != nil && len(attr.Attribute) > 0 {
				printAttributeTable(attr.Attribute)
			}
		}
	}

	return cmd
}

func registerKeyPairCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "key-pair VALUE",
		Aliases: []string{"pair", "kp"},
		Short:   "Register a private and a public key objects from private key PEM encoded data",
		Long: `Register a private and a public key objects from private key PEM encoded data.
		
VALUE can be either plain text, a '-' to read from stdin, or a filename prefixed with @.`,
		Args: cobra.ExactArgs(1),
	}

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

	cmd.Run = func(cmd *cobra.Command, args []string) {
		key := flagsmgmt.BytesFromArg(args[0], 16_000)

		// Register private key
		privReq := kmipClient.Register().PemPrivateKey(key, privateUsage.ToCryptographicUsageMask()).
			WithAttribute(kmip.AttributeNameExtractable, *privateExtractable).
			WithAttribute(kmip.AttributeNameSensitive, *privateSensitive)
		if *privateName != "" {
			privReq = privReq.WithName(*privateName)
		}
		if *description != "" {
			privReq = privReq.WithAttribute(kmip.AttributeNameDescription, *description)
		}
		if *comment != "" {
			privReq = privReq.WithAttribute(kmip.AttributeNameComment, *comment)
		}
		privResp := exit.OnErr2(privReq.ExecContext(cmd.Context()))

		// Register public key
		pubReq := kmipClient.Register().PemPublicKey(key, publicUsage.ToCryptographicUsageMask()).
			WithLink(kmip.LinkTypePrivateKeyLink, privResp.UniqueIdentifier)
		if *publicName != "" {
			pubReq = pubReq.WithName(*publicName)
		}
		if *description != "" {
			pubReq = pubReq.WithAttribute(kmip.AttributeNameDescription, *description)
		}
		if *comment != "" {
			pubReq = pubReq.WithAttribute(kmip.AttributeNameComment, *comment)
		}
		pubResp := exit.OnErr2(pubReq.ExecContext(cmd.Context()))

		// Update public key link in private key
		exit.OnErr2(kmipClient.AddAttribute(privResp.UniqueIdentifier, kmip.AttributeNameLink, kmip.Link{
			LinkType:               kmip.LinkTypePublicKeyLink,
			LinkedObjectIdentifier: pubResp.UniqueIdentifier,
		}).ExecContext(cmd.Context()))

		resp := &payloads.CreateKeyPairResponsePayload{
			PrivateKeyUniqueIdentifier:  privResp.UniqueIdentifier,
			PublicKeyUniqueIdentifier:   pubResp.UniqueIdentifier,
			PrivateKeyTemplateAttribute: privResp.TemplateAttribute,
			PublicKeyTemplateAttribute:  pubResp.TemplateAttribute,
		}

		if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
			output.JsonPrint(resp)
		} else {
			fmt.Println("Pubic Key registered with ID:", resp.PublicKeyUniqueIdentifier)
			fmt.Println("Private Key registered with ID:", resp.PrivateKeyUniqueIdentifier)
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
