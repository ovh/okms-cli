package keys

import (
	"context"
	"encoding/base64"
	"io"

	"github.com/google/uuid"
	"github.com/ovh/okms-cli/cmd/okms/common"
	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/ovh/okms-cli/common/output"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/ovh/okms-sdk-go"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func newEncryptWithServiceKeyCmd() *cobra.Command {
	var (
		useWrap    bool
		noProgress bool
		toBase64   bool
		context    string
	)

	cmd := &cobra.Command{
		Use:   "encrypt KEY-ID DATA [OUTPUT]",
		Short: "Encrypts data, up to 4Kb in size, using provided domain key",
		Long: `Encrypts data, up to 4Kb in size, using provided domain key.

DATA can be either plain text, a '-' to read from stdin, or a filename prefixed with @.
OUTPUT can be either a filepath, or a "-" for stdout. If not set, output is stdout.
`,
		Args: cobra.RangeArgs(2, 3),
		Run: func(cmd *cobra.Command, args []string) {
			keyId := exit.OnErr2(uuid.Parse(args[0]))
			out := "-"
			if len(args) > 2 && args[2] != "" {
				out = args[2]
			}

			if useWrap {
				// If wrapping is enabled, we don't call Encrypt endpoint directly,
				// but instead create a data key, and use it to encrypt data
				ctx := []byte{}
				if context != "" {
					ctx = []byte(context)
				}
				exit.OnErr(wrapEncrypt(cmd.Context(), args[1], out, keyId, ctx, noProgress, toBase64))
				return
			}
			data := flagsmgmt.BytesFromArg(args[1], 8192)
			resp := exit.OnErr2(common.Client().Encrypt(cmd.Context(), common.GetOkmsId(), keyId, context, data))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else {
				writer := flagsmgmt.WriterFromArg(out)
				defer writer.Close()
				exit.OnErr2(writer.Write([]byte(resp)))
			}
		},
	}
	cmd.Flags().BoolVar(&useWrap, "dk", false, "Encrypt locally using a new datakey")
	cmd.Flags().BoolVar(&noProgress, "no-progress", false, "Do not display progress bar or spinner")
	cmd.Flags().BoolVar(&toBase64, "base64", false, "Base64 encode the output when using a datakey")
	cmd.Flags().StringVar(&context, "context", "", "Optional encryption context (AAD)")
	return cmd
}

func wrapEncrypt(ctx context.Context, input, output string, keyId uuid.UUID, keyCtx []byte, noProgress, b64 bool) error {
	in, size := flagsmgmt.ReaderFromArgWithSize(input)
	defer in.Close()
	if !noProgress && output != "-" {
		bar := progressbar.DefaultBytes(size, "Encrypting")
		bReader := progressbar.NewReader(in, bar)
		in = &bReader
	}

	out := flagsmgmt.WriterFromArg(output)
	defer out.Close()
	if b64 {
		out = base64.NewEncoder(base64.StdEncoding, out)
		defer out.Close()
	}

	out, err := common.Client().DataKeys(common.GetOkmsId(), keyId).EncryptStream(ctx, out, keyCtx, okms.BlockSize4MB)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
