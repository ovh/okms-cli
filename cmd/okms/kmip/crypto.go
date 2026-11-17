package kmip

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ovh/kmip-go"
	"github.com/ovh/kmip-go/kmipclient"
	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/ovh/okms-cli/common/flagsmgmt/kmipflags"
	"github.com/ovh/okms-cli/common/output"
	"github.com/ovh/okms-cli/common/utils"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/spf13/cobra"
)

// maxCryptoDataSize is the maximum size of the data read from a CLI argument by the
// encrypt, decrypt, sign and verify commands.
const maxCryptoDataSize = 1 << 20

const cryptoParamsHelp = `Without any of the --mode, --padding, --hash and --tag-length flags, the server uses the
Cryptographic Parameters attribute of the key. Otherwise the parameters given by these flags replace the
key ones: AES keys need --mode (GCM, CBC or ECB, plus --padding for CBC and ECB), RSA keys need --padding
(PKCS1v15, or OAEP with --hash).`

type encryptionFlags struct {
	mode      kmipflags.BlockCipherMode
	padding   kmipflags.PaddingMethod
	hash      kmipflags.HashingAlgorithm
	tagLength int32
	iv        string
	aad       string
}

func setEncryptionFlags(cmd *cobra.Command) *encryptionFlags {
	f := new(encryptionFlags)
	cmd.Flags().Var(&f.mode, "mode", "Block cipher mode, for AES keys")
	cmd.Flags().Var(&f.padding, "padding", "Padding method: None or PKCS5 for AES ECB and CBC modes, PKCS1v15 or OAEP for RSA keys")
	cmd.Flags().Var(&f.hash, "hash", "Hashing algorithm of the OAEP padding")
	cmd.Flags().Int32Var(&f.tagLength, "tag-length", 0, "Length in bytes of the GCM authentication tag")
	cmd.Flags().StringVar(&f.iv, "iv", "", "Base64 encoded IV, counter or nonce")
	cmd.Flags().StringVar(&f.aad, "aad", "", "Additional authenticated data of the GCM mode. Can be either plain text, a '-' to read from stdin, or a filename prefixed with @")
	return f
}

// cryptographicParameters returns the cryptographic parameters of the request, or nil when none of the
// flags selecting them is set, so that the server uses the Cryptographic Parameters attribute of the key.
// When randomIV is true and no IV is given, the server is asked to generate the IV of the GCM and CBC modes.
func (f *encryptionFlags) cryptographicParameters(cmd *cobra.Command, randomIV bool) *kmip.CryptographicParameters {
	changed := false
	for _, name := range []string{"mode", "padding", "hash", "tag-length"} {
		changed = changed || cmd.Flags().Changed(name)
	}
	if !changed {
		return nil
	}

	params := &kmip.CryptographicParameters{
		BlockCipherMode:  kmip.BlockCipherMode(f.mode),
		PaddingMethod:    kmip.PaddingMethod(f.padding),
		HashingAlgorithm: kmip.HashingAlgorithm(f.hash),
		TagLength:        f.tagLength,
	}
	if f.padding == kmipflags.PaddingOAEP {
		if f.hash == 0 {
			exit.OnErr(errors.New("--hash is required with the OAEP padding"))
		}
		params.MaskGenerator = kmip.MaskGeneratorMGF1
		params.MaskGeneratorHashingAlgorithm = params.HashingAlgorithm
	}
	if randomIV && f.iv == "" && (f.mode == kmipflags.GCM || f.mode == kmipflags.CBC) {
		params.RandomIV = utils.PtrTo(true)
	}
	return params
}

// bytesFromArg reads a non empty CLI argument with [flagsmgmt.BytesFromArg].
func bytesFromArg(name, value string) []byte {
	if value == "" {
		exit.Now("%s cannot be empty", name)
	}
	return flagsmgmt.BytesFromArg(value, maxCryptoDataSize)
}

// base64FromArg reads a base64 encoded CLI argument with [flagsmgmt.BytesFromArg] and decodes it.
func base64FromArg(name, value string) []byte {
	raw := strings.TrimSpace(string(bytesFromArg(name, value)))
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		exit.Now("%s is not valid base64: %s", name, err)
	}
	return decoded
}

func encryptCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encrypt ID DATA",
		Short: "Encrypt data with an AES or RSA key",
		Long: `Encrypt data with an AES or RSA key.

DATA can be either plain text, a '-' to read from stdin, or a filename prefixed with @.

` + cryptoParamsHelp + `
For the GCM and CBC modes, the server generates the IV when --iv is not given.

The text output gives the base64 encoded encrypted data, IV and authentication tag, to pass to the decrypt command.`,
		Args: cobra.ExactArgs(2),
	}
	flags := setEncryptionFlags(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		data := bytesFromArg("DATA", args[1])
		req := kmipClient.Encrypt(args[0])
		if params := flags.cryptographicParameters(cmd, true); params != nil {
			req = req.WithCryptographicParameters(*params)
		}
		if flags.iv != "" {
			req = req.WithIvCounterNonce(base64FromArg("--iv", flags.iv))
		}
		if flags.aad != "" {
			req = req.WithAAD(bytesFromArg("--aad", flags.aad))
		}

		resp := exit.OnErr2(req.Data(data).ExecContext(cmd.Context()))
		if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
			output.JsonPrint(resp)
			return
		}
		fmt.Println("Data:", base64.StdEncoding.EncodeToString(resp.Data))
		if len(resp.IVCounterNonce) > 0 {
			fmt.Println("IV:", base64.StdEncoding.EncodeToString(resp.IVCounterNonce))
		}
		if len(resp.AuthenticatedEncryptionTag) > 0 {
			fmt.Println("Tag:", base64.StdEncoding.EncodeToString(resp.AuthenticatedEncryptionTag))
		}
	}

	return cmd
}

func decryptCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decrypt ID DATA",
		Short: "Decrypt data with an AES or RSA key",
		Long: `Decrypt data with an AES or RSA key.

DATA is the base64 encoded encrypted data. It can be either plain text, a '-' to read from stdin, or a filename prefixed with @.

` + cryptoParamsHelp + `
The GCM and CBC modes need the --iv returned by the encrypt command, the GCM mode also needs its --tag.

The text output writes the decrypted data as is.`,
		Args: cobra.ExactArgs(2),
	}
	flags := setEncryptionFlags(cmd)
	tag := cmd.Flags().String("tag", "", "Base64 encoded authentication tag of the GCM mode")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		data := base64FromArg("DATA", args[1])
		req := kmipClient.Decrypt(args[0])
		if params := flags.cryptographicParameters(cmd, false); params != nil {
			req = req.WithCryptographicParameters(*params)
		}
		if flags.iv != "" {
			req = req.WithIvCounterNonce(base64FromArg("--iv", flags.iv))
		}
		if flags.aad != "" {
			req = req.WithAAD(bytesFromArg("--aad", flags.aad))
		}
		if *tag != "" {
			req = req.WithAuthTag(base64FromArg("--tag", *tag))
		}

		resp := exit.OnErr2(req.Data(data).ExecContext(cmd.Context()))
		if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
			output.JsonPrint(resp)
			return
		}
		exit.OnErr2(os.Stdout.Write(resp.Data))
	}

	return cmd
}

const signDataHelp = `DATA can be either plain text, a '-' to read from stdin, or a filename prefixed with @. The server hashes it
with the hash of the signature algorithm. With --digested, DATA is instead the base64 encoded digest of the data,
computed with that same hash.`

func signCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign ID DATA",
		Short: "Sign data with a private key",
		Long: `Sign data with a private key.

` + signDataHelp + `

The text output gives the base64 encoded signature.`,
		Args: cobra.ExactArgs(2),
	}

	var alg kmipflags.SignatureAlgorithm
	cmd.Flags().VarP(&alg, "alg", "a", "Signature algorithm")
	_ = cmd.MarkFlagRequired("alg")
	digested := cmd.Flags().Bool("digested", false, "DATA is a base64 encoded digest")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		req := kmipClient.Sign(args[0]).WithCryptographicParameters(alg.CryptographicParameters())
		var exec kmipclient.ExecSign
		if *digested {
			exec = req.DigestedData(base64FromArg("DATA", args[1]))
		} else {
			exec = req.Data(bytesFromArg("DATA", args[1]))
		}

		resp := exit.OnErr2(exec.ExecContext(cmd.Context()))
		if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
			output.JsonPrint(resp)
			return
		}
		fmt.Println(base64.StdEncoding.EncodeToString(resp.SignatureData))
	}

	return cmd
}

func verifyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "verify ID DATA SIGNATURE",
		Aliases: []string{"signature-verify"},
		Short:   "Verify a signature with a public key",
		Long: `Verify a signature with a public key.

` + signDataHelp + `

SIGNATURE is the base64 encoded signature. It can also be passed from a file or stdin using '-' or '@'.
Stdin can however be only used for 1 argument.

The command fails when the signature is not valid.`,
		Args: cobra.ExactArgs(3),
	}

	var alg kmipflags.SignatureAlgorithm
	cmd.Flags().VarP(&alg, "alg", "a", "Signature algorithm")
	_ = cmd.MarkFlagRequired("alg")
	digested := cmd.Flags().Bool("digested", false, "DATA is a base64 encoded digest")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		req := kmipClient.SignatureVerify(args[0]).WithCryptographicParameters(alg.CryptographicParameters())
		var withData kmipclient.ExecSignatureVerifyWantsSignature
		if *digested {
			withData = req.DigestedData(base64FromArg("DATA", args[1]))
		} else {
			withData = req.Data(bytesFromArg("DATA", args[1]))
		}
		signature := base64FromArg("SIGNATURE", args[2])

		resp := exit.OnErr2(withData.Signature(signature).ExecContext(cmd.Context()))
		jsonOutput := cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT)
		if jsonOutput {
			output.JsonPrint(resp)
		}
		switch resp.ValidityIndicator {
		case kmip.ValidityIndicatorValid:
			if !jsonOutput {
				fmt.Println("Signature is valid")
			}
		case kmip.ValidityIndicatorInvalid:
			exit.OnErr(errors.New("Signature is invalid"))
		default:
			exit.OnErr(errors.New("Signature validity is unknown"))
		}
	}

	return cmd
}
