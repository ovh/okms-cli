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
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func newDecryptWithServiceKeyCmd() *cobra.Command {
	var (
		useWrap    bool
		noProgress bool
		fromBase64 bool
		context    string
	)

	cmd := &cobra.Command{
		Use:   "decrypt KEY-ID DATA [OUTPUT]",
		Short: "Decrypt data previously encrypted by Encrypt operation",
		Long: `Decrypt data previously encrypted by Encrypt operation.

DATA can be either plain text, a '-' to read from stdin, or a filename prefixed with @.
OUTPUT can be either a filepath, or a "-" for stdout. If not set, output is stdout.`,
		Args: cobra.RangeArgs(2, 3),
		Run: func(cmd *cobra.Command, args []string) {
			keyId := exit.OnErr2(uuid.Parse(args[0]))
			out := "-"
			if len(args) > 2 && args[2] != "" {
				out = args[2]
			}

			if useWrap {
				// If wrapping is enabled, we don't call Decrypt endpoint directly,
				// but instead extract data kaye from the blob, and use it to decrypt data
				ctx := []byte{}
				if context != "" {
					ctx = []byte(context)
				}
				exit.OnErr(wrapDecrypt(cmd.Context(), args[1], out, keyId, ctx, noProgress, fromBase64))
				return
			}

			text := flagsmgmt.BytesFromArg(args[1], 8192)
			resp := exit.OnErr2(common.Client().Decrypt(cmd.Context(), common.GetOkmsId(), keyId, context, string(text)))
			if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
				output.JsonPrint(resp)
			} else {
				writer := flagsmgmt.WriterFromArg(out)
				defer writer.Close()
				exit.OnErr2(writer.Write(resp))
			}
		},
	}

	cmd.Flags().BoolVar(&useWrap, "dk", false, "Decrypt locally using an embedded encrypted datakey")
	cmd.Flags().BoolVar(&noProgress, "no-progress", false, "Do not display progress bar or spinner")
	cmd.Flags().BoolVar(&fromBase64, "base64", false, "When using a datakey, decrypts a base64 encoded input")
	cmd.Flags().StringVar(&context, "context", "", "Optional encryption context (AAD)")
	return cmd
}

func wrapDecrypt(ctx context.Context, input, output string, keyId uuid.UUID, keyCtx []byte, noProgress, b64 bool) error {
	reader, size := flagsmgmt.ReaderFromArgWithSize(input)
	defer reader.Close()
	if !noProgress && output != "-" {
		bar := progressbar.DefaultBytes(size, "Decrypting")
		bReader := progressbar.NewReader(reader, bar)
		reader = &bReader
	}

	var in io.Reader = reader
	if b64 {
		in = base64.NewDecoder(base64.StdEncoding, reader)
	}

	out := flagsmgmt.WriterFromArg(output)
	defer out.Close()

	in, err := common.Client().DataKeys(common.GetOkmsId(), keyId).DecryptStream(ctx, in, keyCtx)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, in)
	return err
}
