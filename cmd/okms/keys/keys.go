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
		keysPageSize uint32
		listAll      bool
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
			for key, err := range common.Client().ListAllServiceKeys(&keysPageSize, &stateFilter).Iter(cmd.Context()) {
				exit.OnErr(err)
				keys.ObjectsList = append(keys.ObjectsList, *key)
			}

			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(keys)
			} else {
				table := tablewriter.NewWriter(os.Stdout)
				table.Header([]string{"ID", "Name", "Type", "State", "Created At"})
				for _, key := range keys.ObjectsList {
					keyAttr := getCommonKeyAttributes(&key)
					exit.OnErr(table.Append([]string{
						key.Id.String(),
						key.Name,
						string(key.Type),
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

	cmd.Flags().Uint32Var(&keysPageSize, "page-size", 100, "Number of keys to fetch per page (between 10 and 500)")
	cmd.Flags().BoolVarP(&listAll, "all", "A", false, "List all keys (including deactivated and deleted ones)")
	return cmd
}

func newAddServiceKeyCmd() *cobra.Command {
	var (
		keyUsage restflags.KeyUsageList
		keySize  int32
		//lint:ignore ST1023 setting default
		keySpec    = restflags.OCTETSTREAM
		curveType  restflags.CurveType
		keyContext string
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
			if keySpec == restflags.ELLIPTIC_CURVE {
				crv := curveType.ToRestCurve()
				body.Curve = &crv
			} else {
				keySizeEnum := types.KeySizes(keySize)
				body.Size = &keySizeEnum
			}

			resp := exit.OnErr2(common.Client().CreateImportServiceKey(cmd.Context(), nil, body))
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
	cmd.MarkFlagsMutuallyExclusive("size", "curve")
	return cmd
}

func newGetServiceKeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get KEY-ID",
		Short: "Retrieve domain key metadata",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			keyId := exit.OnErr2(uuid.Parse(args[0]))
			resp := exit.OnErr2(common.Client().GetServiceKey(cmd.Context(), keyId, nil))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else {
				printServiceKey(resp)
			}
		},
	}
}

func newExportPublicKeyCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "export KEY-ID",
		Short: "Export public key material",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			keyId := exit.OnErr2(uuid.Parse(args[0]))
			resp := exit.OnErr2(common.Client().GetServiceKey(cmd.Context(), keyId, utils.PtrTo(types.Jwk)))
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
		keyUsage   restflags.KeyUsageList
		symmetric  bool
		keyContext string
	)
	cmd := &cobra.Command{
		Use:   "import NAME KEY",
		Short: "Import a private base64 encoded symmetric key or a PEM encoded assymmetric key",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if keyContext == "" {
				keyContext = args[0]
			}
			key := flagsmgmt.BytesFromArg(args[1], 8192)
			var resp *types.GetServiceKeyResponse
			if !symmetric {
				resp = exit.OnErr2(common.Client().ImportKeyPairPEM(cmd.Context(), key, args[0], keyContext, keyUsage.ToCryptographicUsage()...))
			} else {
				k := exit.OnErr2(base64.StdEncoding.DecodeString(string(key)))
				resp = exit.OnErr2(common.Client().ImportKey(cmd.Context(), k, args[0], keyContext, keyUsage.ToCryptographicUsage()...))
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
		{"Key Type", string(kt)},
		{"Size", size},
		{"Curve", curve},
		{"Usage", usage},
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
					if err := common.Client().DeactivateServiceKey(cmd.Context(), keyId, types.Unspecified); err != nil {
						errs = append(errs, fmt.Errorf("Failed to deactivate key %q: %w", id, err))
						continue
					}
				}
				if err := common.Client().DeleteServiceKey(cmd.Context(), keyId); err != nil {
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
				if err := common.Client().DeactivateServiceKey(cmd.Context(), keyId, revocationReason.RestModel()); err != nil {
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
				if err := common.Client().ActivateServiceKey(cmd.Context(), keyId); err != nil {
					errs = append(errs, fmt.Errorf("Failed to activate key %q: %w", id, err))
				}
			}
			exit.OnErr(errors.Join(errs...))
		},
	}
}

func newUpdateKeyCmd() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "update KEY-ID",
		Args:  cobra.ExactArgs(1),
		Short: "Update a service key",
		Run: func(cmd *cobra.Command, args []string) {
			keyId := exit.OnErr2(uuid.Parse(args[0]))
			body := types.PatchServiceKeyRequest{}
			if name != "" {
				body.Name = name
			}
			resp := exit.OnErr2(common.Client().UpdateServiceKey(cmd.Context(), keyId, body))
			printServiceKey(resp)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Update key with a new name")
	if err := cmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	return cmd
}

type KeyAttr struct {
	CreatedAt     time.Time
	State         types.KeyStates
	ActivatedAt   *time.Time
	CompromisedAt *time.Time
	DeactivatedAt *time.Time
}

func getCommonKeyAttributes(key *types.GetServiceKeyResponse) KeyAttr {
	keyAttr := KeyAttr{}
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
