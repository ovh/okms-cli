package keys

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/olekukonko/tablewriter"
	"github.com/ovh/okms-sdk-go"
	"golang.org/x/crypto/ssh"

	"github.com/ovh/okms-cli/cmd/okms/common"
	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/ovh/okms-cli/common/flagsmgmt/restflags"
	"github.com/ovh/okms-cli/common/output"
	"github.com/ovh/okms-cli/common/utils"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/ovh/okms-sdk-go/types"
	"github.com/spf13/cobra"
)

func newListServiceKeysCmd() *cobra.Command {
	var (
		pageSize uint32
		listAll  bool
	)

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List domain keys",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			keys := types.ListServiceKeysResponse{
				ObjectsList: []types.GetServiceKeyResponse{},
			}
			// Let's list all the keys by putting them all in memory. The memory is not an issue, unless a domain has hundreds of thousands of keys
			// Filter keys by activation state
			stateFilter := types.KeyStatesActive
			if listAll {
				stateFilter = types.KeyStatesAll
			}
			for key, err := range common.Client().ListAllServiceKeys(common.GetOkmsId(), &pageSize, &stateFilter).Iter(cmd.Context()) {
				exit.OnErr(err)
				keys.ObjectsList = append(keys.ObjectsList, *key)
			}

			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(keys)
			} else {
				table := tablewriter.NewWriter(os.Stdout)
				table.Header([]string{"ID", "Name", "Type", "Class", "State", "Created At"})
				for _, key := range keys.ObjectsList {
					keyAttr := getCommonKeyAttributes(&key)
					exit.OnErr(table.Append([]string{
						key.Id.String(),
						key.Name,
						string(key.Type),
						string(keyAttr.Class),
						string(keyAttr.State),
						keyAttr.CreatedAt.Format(time.DateTime),
						// strconv.Itoa(int(*key.KeySize)),
						// strings.Join(*key.KeyOps, ", "),
						// strconv.Itoa(int(*key.LatestVersion)),
					}))
				}
				exit.OnErr(table.Render())
			}
		},
	}

	cmd.Flags().Uint32Var(&pageSize, "page-size", 100, "Number of keys to fetch per page (between 10 and 500)")
	cmd.Flags().BoolVarP(&listAll, "all", "A", false, "List all keys (including deactivated and deleted ones)")
	return cmd
}

func newAddServiceKeyCmd() *cobra.Command {
	var (
		keyUsage restflags.KeyUsageList
		keySize  int32
		//lint:ignore ST1023 setting default
		keySpec         = restflags.OCTETSTREAM
		curveType       restflags.CurveType
		protectionLevel restflags.ProtectionLevel
		keyContext      string
		keyID           string
		extractable     bool
	)

	cmd := &cobra.Command{
		Use:     "generate NAME",
		Short:   "Generate a new domain service key",
		Aliases: []string{"new", "gen", "create", "add"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if keyContext == "" {
				// Use the key name as the default context.
				keyContext = args[0]
			}
			body := types.CreateImportServiceKeyRequest{
				Context:    &keyContext,
				Name:       args[0],
				Type:       utils.PtrTo(keySpec.RestModel()),
				Operations: utils.PtrTo(keyUsage.ToCryptographicUsage()),
				Keys:       nil,
			}
			if protectionLevel != "" {
				body.ProtectionLevel = (*types.ProtectionLevelEnum)(&protectionLevel)
			}

			if keySpec == restflags.ELLIPTIC_CURVE {
				crv := curveType.ToRestCurve()
				body.Curve = &crv
			} else {
				keySizeEnum := types.KeySizes(keySize)
				body.Size = &keySizeEnum
			}

			if keyID != "" {
				id := exit.OnErr2(uuid.Parse(keyID))
				body.Id = utils.PtrTo(id.String())
			}

			body.Extractable = utils.PtrTo(extractable)

			resp := exit.OnErr2(common.Client().CreateImportServiceKey(cmd.Context(), common.GetOkmsId(), nil, body))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else {
				printServiceKey(resp)
			}
		},
	}

	cmd.Flags().StringVar(&keyContext, "context", "", "Context of the key. Defaults to the key's name")
	cmd.Flags().Var(&keyUsage, "usage", "Key operations (Key usage).")

	cmd.Flags().Var(&keySpec, "type", "Defines type of a key to be created.")
	cmd.Flags().Int32Var(&keySize, "size", 256, "Size of the key to be generated")
	cmd.Flags().Var(&curveType, "curve", "Curve type for Elliptic Curve (ec) keys.")
	cmd.Flags().Var(&protectionLevel, "protectionLevel", "Level of protection of the key's storage (software or HSM).")
	cmd.Flags().StringVar(&keyID, "keyId", "", "Optional key ID (UUID)")
	cmd.Flags().BoolVar(&extractable, "extractable", false, "Whether the key and its material can be extracted (exported plain or wrapped). Defaults to false.")
	cmd.MarkFlagsMutuallyExclusive("size", "curve")
	return cmd
}

func newGetServiceKeyCmd() *cobra.Command {
	var (
		wrappingKeyID     string
		wrappingAlgorithm string
		wrappedKeyFormat  string
	)

	cmd := &cobra.Command{
		Use:   "get KEY-ID",
		Short: "Retrieve domain key metadata, or export the key material in wrapped form",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			keyId := exit.OnErr2(uuid.Parse(args[0]))

			// When a wrapping key is provided, export the key material in wrapped (encrypted) form.
			if wrappingKeyID != "" {
				wrapKeyId := exit.OnErr2(uuid.Parse(wrappingKeyID))

				algo := types.WrappingAlgorithms(wrappingAlgorithm)
				if !algo.Valid() {
					exit.OnErr(fmt.Errorf("Invalid wrapping algorithm %q, expected one of [RSA-OAEP|RSA-OAEP-256]", wrappingAlgorithm))
				}
				format := types.KeyFormatTypes(wrappedKeyFormat)
				if !format.Valid() {
					exit.OnErr(fmt.Errorf("Invalid wrapped key format %q, expected one of [JWK|RAW|PKCS1|PKCS8]", wrappedKeyFormat))
				}

				wrappedKeys := exit.OnErr2(common.Client().GetWrappedServiceKey(cmd.Context(), common.GetOkmsId(), keyId, wrapKeyId, format, algo))
				if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
					output.JsonPrint(wrappedKeys)
				} else {
					for _, k := range wrappedKeys {
						fmt.Println(k.Ciphertext)
					}
				}
				return
			}

			resp := exit.OnErr2(common.Client().GetServiceKey(cmd.Context(), common.GetOkmsId(), keyId, nil))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else {
				printServiceKey(resp)
			}
		},
	}

	cmd.Flags().StringVar(&wrappingKeyID, "wrapping-key-id", "", "ID of the transport (wrapping) key used to encrypt the exported key material. When set, the key is exported in wrapped form instead of returning metadata")
	cmd.Flags().StringVar(&wrappingAlgorithm, "wrapping-algorithm", string(types.RSAOAEP256), "Key wrapping algorithm [RSA-OAEP|RSA-OAEP-256]")
	cmd.Flags().StringVar(&wrappedKeyFormat, "wrapped-key-format", string(types.JWK), "Format of the plaintext key material before wrapping [JWK|RAW|PKCS1|PKCS8]")

	return cmd
}

func newExportPublicKeyCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "export KEY-ID",
		Short: "Export public key material",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			keyId := exit.OnErr2(uuid.Parse(args[0]))
			resp := exit.OnErr2(common.Client().GetServiceKey(cmd.Context(), common.GetOkmsId(), keyId, utils.PtrTo(types.Jwk)))
			if resp.Keys == nil || len(*resp.Keys) == 0 {
				exit.OnErr(errors.New("Server returned no key"))
			}
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
				return
			}

			if resp.Attributes != nil && (*resp.Attributes)["state"] != "active" {
				exit.OnErr(fmt.Errorf("The key is not active (state is %q)", (*resp.Attributes)["state"]))
			}

			key := (*resp.Keys)[0]

			if strings.EqualFold(format, "jwk") {
				output.JsonPrint(key)
				return
			}
			rawKey := exit.OnErr2(key.PublicKey())

			if strings.EqualFold(format, "pkcs1") {
				if rsaKey, ok := rawKey.(*rsa.PublicKey); ok {
					pemBlock := pem.Block{
						Type:  "RSA PUBLIC KEY",
						Bytes: x509.MarshalPKCS1PublicKey(rsaKey),
					}
					exit.OnErr(pem.Encode(os.Stdout, &pemBlock))
					return
				}
				exit.OnErr(errors.New("pkcs1 format is only for RSA public keys"))
			} else if strings.EqualFold(format, "openssh") {
				sshKey := exit.OnErr2(ssh.NewPublicKey(rawKey))
				rawSshKey := bytes.TrimSpace(ssh.MarshalAuthorizedKey(sshKey))
				rawSshKey = append(rawSshKey, append([]byte{' '}, []byte(resp.Name)...)...)
				fmt.Println(string(rawSshKey))
				return
			}

			pemBlock := pem.Block{
				Type:  "PUBLIC KEY",
				Bytes: exit.OnErr2(x509.MarshalPKIXPublicKey(rawKey)),
			}
			exit.OnErr(pem.Encode(os.Stdout, &pemBlock))
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "pkix", "Export format [pkix|pkcs1|openssh|jwk]")

	return cmd
}

func newImportServiceKeyCmd() *cobra.Command {
	var (
		keyUsage         restflags.KeyUsageList
		symmetric        bool
		keyContext       string
		keyID            string
		wrappingKeyID    string
		wrappedKeyFormat string
		extractable      bool
	)
	cmd := &cobra.Command{
		Use:   "import NAME KEY",
		Short: "Import a symmetric, asymmetric (private or public), or wrapped key",
		Long: `Import a key into the domain.

There are two distinct import modes: plain import and wrapped import.

Plain import (default):
  The key material is provided in the clear. KEY is either:
    - a PEM encoded asymmetric key. Supported PEM types are PKCS8, PKCS1, PKIX,
      SEC1 and OpenSSH private keys, as well as PKIX ("PUBLIC KEY"), PKCS1
      ("RSA PUBLIC KEY") and EC public keys. A public-only key is imported
      without private material and can only be used for public operations
      (e.g. verify or encrypt).
    - a base64 encoded symmetric key, when --symmetric is set.

Wrapped import (--wrapping-key-id):
  The key material never appears in the clear. KEY is wrapped (encrypted) key
  material as a JWE Compact Serialization string, produced by wrapping the
  plaintext key with the public part of the transport key identified by
  --wrapping-key-id. --wrapped-key-format describes the format of the plaintext
  key that was wrapped [JWK|RAW|PKCS1|PKCS8]. The KMS unwraps the material and
  infers the key type, size and curve from it.

In both modes KEY may be given inline, as @path to read from a file, or - to
read from stdin.`,
		Example: `  # Plain import of a PEM encoded RSA private key
  okms keys import --usage sign,verify my-rsa @private.pem

  # Plain import of a base64 encoded symmetric key
  okms keys import --symmetric --usage encrypt,decrypt my-aes <base64-key>

  # Plain import of a public-only key (verify/encrypt operations only)
  okms keys import --usage verify my-rsa-pub @public.pem

  # Wrapped import of PKCS8 key material wrapped with transport key <transport-id>
  okms keys import --usage encrypt,decrypt --wrapping-key-id <transport-id> \
      --wrapped-key-format PKCS8 my-imported @wrapped.jwe`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if keyContext == "" {
				keyContext = args[0]
			}
			key := flagsmgmt.BytesFromArg(args[1], 8192)

			var opts []okms.ServiceKeyOption
			if keyID != "" {
				id := exit.OnErr2(uuid.Parse(keyID))
				opts = append(opts, okms.WithKeyID(id))
			}
			opts = append(opts, okms.WithExtractable(extractable))

			var resp *types.GetServiceKeyResponse
			switch {
			case wrappingKeyID != "":
				// KEY is wrapped (encrypted) key material as a JWE Compact Serialization string.
				wrapKeyId := exit.OnErr2(uuid.Parse(wrappingKeyID))
				format := types.KeyFormatTypes(wrappedKeyFormat)
				if !format.Valid() {
					exit.OnErr(fmt.Errorf("Invalid wrapped key format %q, expected one of [JWK|RAW|PKCS1|PKCS8]", wrappedKeyFormat))
				}
				ciphertext := strings.TrimSpace(string(key))
				resp = exit.OnErr2(common.Client().ImportWrappedServiceKey(cmd.Context(), common.GetOkmsId(), wrapKeyId, ciphertext, format, args[0], keyContext, keyUsage.ToCryptographicUsage(), opts...))
			case !symmetric:
				resp = exit.OnErr2(common.Client().ImportKeyPairPEM(cmd.Context(), common.GetOkmsId(), key, args[0], keyContext, keyUsage.ToCryptographicUsage(), opts...))
			default:
				k := exit.OnErr2(base64.StdEncoding.DecodeString(string(key)))
				resp = exit.OnErr2(common.Client().ImportKey(cmd.Context(), common.GetOkmsId(), k, args[0], keyContext, keyUsage.ToCryptographicUsage(), opts...))
			}

			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else {
				printServiceKey(resp)
			}
		},
	}

	cmd.Flags().StringVar(&keyContext, "context", "", "Context of the key. Defaults to the key's name")
	cmd.Flags().Var(&keyUsage, "usage", "Key operations (Key usage).")
	cmd.Flags().BoolVarP(&symmetric, "symmetric", "S", false, "Import a base64 encoded symmetric key")
	cmd.Flags().StringVar(&wrappingKeyID, "wrapping-key-id", "", "ID of the transport (wrapping) key that was used to wrap KEY. When set, KEY is imported as wrapped (JWE) key material and its type is inferred by the KMS")
	cmd.Flags().StringVar(&wrappedKeyFormat, "wrapped-key-format", string(types.JWK), "Format of the plaintext key material that was wrapped [JWK|RAW|PKCS1|PKCS8]")
	cmd.Flags().BoolVar(&extractable, "extractable", false, "Whether the imported key and its material can be extracted (exported plain or wrapped). Defaults to false.")
	cmd.MarkFlagsMutuallyExclusive("symmetric", "wrapping-key-id")
	return cmd
}

func printServiceKey(resp *types.GetServiceKeyResponse) {
	id := resp.Id
	name := resp.Name
	keyAttr := getCommonKeyAttributes(resp)
	kt := resp.Type
	var size string
	if resp.Size != nil {
		size = fmt.Sprintf("%d", *resp.Size)
	}
	var curve string
	if resp.Curve != nil {
		curve = string(*resp.Curve)
	}

	var usage string
	if resp.Operations != nil {
		ops := make([]string, len(*resp.Operations))
		for i := range *resp.Operations {
			ops[i] = string((*resp.Operations)[i])
		}
		usage = strings.Join(ops, ", ")
	}

	table := tablewriter.NewWriter(os.Stdout)
	exit.OnErr(table.Bulk([][]string{
		{"Id", id.String()},
		{"Name", name},
		{"State", string(keyAttr.State)},
		{"Class", string(keyAttr.Class)},
		{"Key Type", string(kt)},
		{"Size", size},
		{"Curve", curve},
		{"Usage", usage},
		{"Protection Level", string(resp.ProtectionLevel)},
		{"Created at", keyAttr.CreatedAt.Format(time.DateTime)},
	}))
	if keyAttr.ActivatedAt != nil {
		exit.OnErr(table.Append([]string{"Activated at", keyAttr.ActivatedAt.Format(time.DateTime)}))
	}
	if keyAttr.DeactivatedAt != nil {
		exit.OnErr(table.Append([]string{"Deactivated at", keyAttr.DeactivatedAt.Format(time.DateTime)}))
	}
	if keyAttr.CompromisedAt != nil {
		exit.OnErr(table.Append([]string{"Compromised at", keyAttr.CompromisedAt.Format(time.DateTime)}))
	}
	exit.OnErr(table.Render())
}

func newDeleteKeyCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:     "delete KEY-ID [KEY-ID...]",
		Aliases: []string{"del"},
		Args:    cobra.MinimumNArgs(1),
		Short:   "Delete one or more deactivated service keys. This action is irreversible",
		Run: func(cmd *cobra.Command, args []string) {
			var errs []error
			for _, id := range args {
				keyId, err := uuid.Parse(id)
				if err != nil {
					errs = append(errs, fmt.Errorf("Invalid Key ID %q: %w", id, err))
					continue
				}
				if force {
					if err := common.Client().DeactivateServiceKey(cmd.Context(), common.GetOkmsId(), keyId, types.Unspecified); err != nil {
						errs = append(errs, fmt.Errorf("Failed to deactivate key %q: %w", id, err))
						continue
					}
				}
				if err := common.Client().DeleteServiceKey(cmd.Context(), common.GetOkmsId(), keyId); err != nil {
					errs = append(errs, fmt.Errorf("Failed to delete key %q: %w", id, err))
				}
			}
			exit.OnErr(errors.Join(errs...))
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Force delete on active keys by deactivating them first with an unspecified reason")

	return cmd
}

func newDeactivateKeyCmd() *cobra.Command {
	//lint:ignore ST1023 for readability
	var revocationReason = restflags.Unspecified
	cmd := &cobra.Command{
		Use:   "deactivate KEY-ID [KEY-ID...]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Deactivate one or more service keys",
		Run: func(cmd *cobra.Command, args []string) {
			var errs []error
			for _, id := range args {
				keyId, err := uuid.Parse(id)
				if err != nil {
					errs = append(errs, fmt.Errorf("Invalid Key ID %q: %w", id, err))
					continue
				}
				if err := common.Client().DeactivateServiceKey(cmd.Context(), common.GetOkmsId(), keyId, revocationReason.RestModel()); err != nil {
					errs = append(errs, fmt.Errorf("Failed to deactivate key %q: %w", id, err))
				}
			}
			exit.OnErr(errors.Join(errs...))
		},
	}
	cmd.Flags().Var(&revocationReason, "reason", "The reason of revocation")
	return cmd
}

func newActivateKeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "activate KEY-ID [KEY-ID...]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Activate one or more service keys",
		Run: func(cmd *cobra.Command, args []string) {
			var errs []error
			for _, id := range args {
				keyId, err := uuid.Parse(id)
				if err != nil {
					errs = append(errs, fmt.Errorf("Invalid Key ID %q: %w", id, err))
					continue
				}
				if err := common.Client().ActivateServiceKey(cmd.Context(), common.GetOkmsId(), keyId); err != nil {
					errs = append(errs, fmt.Errorf("Failed to activate key %q: %w", id, err))
				}
			}
			exit.OnErr(errors.Join(errs...))
		},
	}
}

func newUpdateKeyCmd() *cobra.Command {
	var (
		name        string
		extractable bool
	)

	cmd := &cobra.Command{
		Use:   "update KEY-ID",
		Args:  cobra.ExactArgs(1),
		Short: "Update a service key",
		Run: func(cmd *cobra.Command, args []string) {
			keyId := exit.OnErr2(uuid.Parse(args[0]))
			body := types.PatchServiceKeyRequest{}
			if cmd.Flags().Changed("name") {
				body.Name = utils.PtrTo(name)
			}
			if cmd.Flags().Changed("extractable") {
				body.Extractable = utils.PtrTo(extractable)
			}
			resp := exit.OnErr2(common.Client().UpdateServiceKey(cmd.Context(), common.GetOkmsId(), keyId, body))
			printServiceKey(resp)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Update key with a new name")
	cmd.Flags().BoolVar(&extractable, "extractable", false, "Set whether the key and its material can be extracted. Set to false to block plain or wrapped key export")
	cmd.MarkFlagsOneRequired("name", "extractable")

	return cmd
}

type KeyAttr struct {
	CreatedAt     time.Time
	State         types.KeyStates
	Class         types.ServiceKeyClassEnum
	ActivatedAt   *time.Time
	CompromisedAt *time.Time
	DeactivatedAt *time.Time
}

func getCommonKeyAttributes(key *types.GetServiceKeyResponse) KeyAttr {
	keyAttr := KeyAttr{}
	if key.Class != nil {
		keyAttr.Class = *key.Class
	}
	if key.Attributes != nil && *key.Attributes != nil {
		if str, ok := (*key.Attributes)["original_creation_date"].(string); ok {
			keyAttr.CreatedAt, _ = time.Parse(time.RFC3339, str)
		}
		if str, ok := (*key.Attributes)["activation_date"].(string); ok {
			if tm, err := time.Parse(time.RFC3339, str); err == nil {
				keyAttr.ActivatedAt = &tm
			}
		}
		if str, ok := (*key.Attributes)["compromise_date"].(string); ok {
			if tm, err := time.Parse(time.RFC3339, str); err == nil {
				keyAttr.CompromisedAt = &tm
			}
		}
		if str, ok := (*key.Attributes)["deactivation_date"].(string); ok {
			if tm, err := time.Parse(time.RFC3339, str); err == nil {
				keyAttr.DeactivatedAt = &tm
			}
		}
		if state, ok := (*key.Attributes)["state"].(string); ok {
			keyAttr.State = types.KeyStates(state)
		}
	}
	return keyAttr
}
