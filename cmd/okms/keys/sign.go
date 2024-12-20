package keys

import (
	"crypto/ecdsa"
	"crypto/rsa"

	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/google/uuid"
	"github.com/ovh/okms-cli/cmd/okms/common"
	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/ovh/okms-cli/common/flagsmgmt/restflags"
	"github.com/ovh/okms-cli/common/output"
	"github.com/ovh/okms-cli/common/utils"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/ovh/okms-sdk-go/types"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func newSignCmd() *cobra.Command {
	var noProgress bool

	signCmd := &cobra.Command{
		Use:   "sign KEY-ID DATA",
		Args:  cobra.ExactArgs(2),
		Short: "Sign a raw data or a base64 encoded digest with the given key",
		Long: `Sign a raw data or a base64 encoded digest with the given key.

When --digest is unset, DATA must be a base64 encoded digest. But if --digest is given,
then DATA will be hashed using the provided alogorithm.
In both cases, DATA can be either plain text, a '-' to read from stdin, or a filename prefixed with @`,
	}

	params := setSignVerifyCommonFlags(signCmd)

	signCmd.Run = func(cmd *cobra.Command, args []string) {
		data := readDigest(params.signatureAlgorithm, args[1], "Signing", noProgress)
		keyId := exit.OnErr2(uuid.Parse(args[0]))
		signature := exit.OnErr2(common.Client().Sign(cmd.Context(), keyId, nil, params.signatureAlgorithm.Alg(), true, data))
		if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
			output.JsonPrint(signature)
		} else {
			fmt.Println(signature)
		}
	}

	signCmd.Flags().BoolVar(&noProgress, "no-progress", false, "Do not display progress bar or spinner")

	return signCmd
}

func newVerifyCmd() *cobra.Command {
	var (
		noProgress bool
		local      bool
	)

	verifyCmd := &cobra.Command{
		Use:   "verify KEY-ID DATA SIGNATURE",
		Args:  cobra.ExactArgs(3),
		Short: "Verify a signature against a key and a raw data or a base64 encoded digest",
		Long: `Verify a signature against a key and a raw data or a base64 encoded digest.

When --digest is unset, DATA must be a base64 encoded digest. But if --digest is given,
then DATA will be hashed using the provided alogorithm.
In both cases, DATA can be either plain text, a '-' to read from stdin, or a filename prefixed with @.

SIGNATURE can also be passed from a file or stdin using '-' or '@'. Stdin can however be only used for 1 argument. 
`,
	}

	params := setSignVerifyCommonFlags(verifyCmd)

	verifyCmd.Run = func(cmd *cobra.Command, args []string) {
		data := readDigest(params.signatureAlgorithm, args[1], "Verifying signature", noProgress)
		signature := flagsmgmt.StringFromArg(args[2], 8192)
		keyId := exit.OnErr2(uuid.Parse(args[0]))
		if !local {
			valid := exit.OnErr2(common.Client().Verify(cmd.Context(), keyId, params.signatureAlgorithm.Alg(), true, data, signature))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(valid)
			} else if valid {
				fmt.Println("Signature is valid")
			}
			if !valid {
				exit.OnErr(errors.New("Signature invalid"))
			}
		} else {
			resp := exit.OnErr2(common.Client().GetServiceKey(cmd.Context(), keyId, utils.PtrTo(types.Jwk)))
			if len(*resp.Keys) == 0 {
				exit.OnErr(errors.New("The server returned no key"))
			}
			key := (*resp.Keys)[0]

			hashAlg := params.signatureAlgorithm.HashAlgorithm()

			sig := exit.OnErr2(base64.StdEncoding.DecodeString(signature))

			rawKey := exit.OnErr2(key.PublicKey())
			switch k := rawKey.(type) {
			case *rsa.PublicKey:
				var err error
				switch params.signatureAlgorithm {
				case restflags.RS256, restflags.RS384, restflags.RS512:
					err = rsa.VerifyPKCS1v15(k, hashAlg, data, sig)
				case restflags.PS256, restflags.PS384, restflags.PS512:
					err = rsa.VerifyPSS(k, hashAlg, data, sig, nil)
				default:
					exit.OnErr(fmt.Errorf("Cannot use algorithm %q with an RSA key", params.signatureAlgorithm))
				}
				if err != nil {
					exit.OnErr(fmt.Errorf("Validation failed: %w", err))
				}
			case *ecdsa.PublicKey:
				// Signature is encoded with format IEEE P1363 instead of ASN.1 (RFC 3279)
				// We need to split it in halves to retrieve R and S parts.
				r := new(big.Int).SetBytes(sig[0 : len(sig)/2])
				s := new(big.Int).SetBytes(sig[len(sig)/2:])
				switch params.signatureAlgorithm {
				case restflags.ES256, restflags.ES384, restflags.ES512:
					if !ecdsa.Verify(k, data, r, s) {
						exit.OnErr(fmt.Errorf("signature is not valid"))
					}
				default:
					exit.OnErr(fmt.Errorf("Cannot use algorithm %q with an EC key", params.signatureAlgorithm))
				}
			default:
				exit.OnErr(fmt.Errorf("unsuported key type %T", rawKey))
			}
			fmt.Println("Signature is valid")
		}
	}

	verifyCmd.Flags().BoolVar(&noProgress, "no-progress", false, "Do not display progress bar or spinner")
	verifyCmd.Flags().BoolVar(&local, "local", false, "Verify the signature localy using the key material")

	return verifyCmd
}

type signVerifyParams struct {
	signatureAlgorithm restflags.SignatureAlgorithm
}

func setSignVerifyCommonFlags(cmd *cobra.Command) *signVerifyParams {
	params := new(signVerifyParams)
	cmd.Flags().VarP(&params.signatureAlgorithm, "alg", "a", "Signature algorithm")
	if err := cmd.MarkFlagRequired("alg"); err != nil {
		panic(err)
	}
	return params
}

func readDigest(digestAlgorithm restflags.SignatureAlgorithm, input, msg string, noProgress bool) []byte {
	d := digestAlgorithm.NewHasher()
	reader, size := flagsmgmt.ReaderFromArgWithSize(input)
	defer reader.Close()
	var writer io.Writer = d
	if !noProgress {
		bar := progressbar.DefaultBytes(size, msg)
		defer bar.Close()
		writer = io.MultiWriter(writer, bar)
	}
	exit.OnErr2(io.Copy(writer, reader))
	return d.Sum(nil)
}
