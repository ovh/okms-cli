package kmip

import (
	"fmt"
	"os"

	"github.com/ovh/kmip-go"
	"github.com/ovh/kmip-go/ttlv"
	"github.com/ovh/okms-cli/common/flagsmgmt"
	"github.com/ovh/okms-cli/common/output"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/spf13/cobra"
)

func getCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get ID",
		Short: "Get the materials from a kmip object",
		Args:  cobra.ExactArgs(1),
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		req := kmipClient.Get(args[0])

		resp := exit.OnErr2(req.ExecContext(cmd.Context()))
		if cmd.Flag("output").Value.String() == string(flagsmgmt.JSON_OUTPUT_FORMAT) {
			output.JsonPrint(resp)
			return
		}

		switch obj := resp.Object.(type) {
		case *kmip.SecretData:
			secret := exit.OnErr2(obj.Data())
			os.Stdout.Write(secret)
		case *kmip.SymmetricKey:
			key := exit.OnErr2(obj.KeyMaterial())
			os.Stdout.Write(key)
		case *kmip.Certificate:
			cert := exit.OnErr2(obj.PemCertificate())
			fmt.Println(cert)
		case *kmip.PrivateKey:
			pem := exit.OnErr2(obj.Pkcs8Pem())
			fmt.Println(pem)
		case *kmip.PublicKey:
			pem := exit.OnErr2(obj.PkixPem())
			fmt.Println(pem)
		default:
			os.Stdout.Write(ttlv.MarshalText(obj))
		}
	}

	return cmd
}
