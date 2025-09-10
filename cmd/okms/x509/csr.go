package x509

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"net"
	"net/url"
	"os"

	"github.com/google/uuid"
	"github.com/ovh/okms-cli/cmd/okms/common"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/spf13/cobra"
)

func createGenerateCsrCommand() *cobra.Command {
	var (
		subject  *pkixNameParams
		dnsNames []string
		emails   []string
		ips      []net.IP
	)

	cmd := &cobra.Command{
		Use:   "csr KEY-ID",
		Short: "Generate a CSR signed with the given private key",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			keyId := exit.OnErr2(uuid.Parse(args[0]))
			signer := exit.OnErr2(common.Client().NewSigner(cmd.Context(), common.GetOkmsId(), keyId))
			template := &x509.CertificateRequest{
				// Do not specify signature algorithm here, and let go chose one based on the key type.
				// SignatureAlgorithm: x509.ECDSAWithSHA384,
				Subject:        subject.pkixName(),
				DNSNames:       dnsNames,
				EmailAddresses: emails,
				IPAddresses:    ips,
				URIs:           []*url.URL{}, //TODO: Make it a configurable option ?
				ExtraExtensions: []pkix.Extension{
					// Embbed subject key identifier in the csr extensions
					{Id: OID_CE_SUBJECT_KEY_IDENTIFIER, Critical: false, Value: exit.OnErr2(asn1.Marshal(keyId[:]))},
				},
			}
			csrBytes := exit.OnErr2(x509.CreateCertificateRequest(rand.Reader, template, signer))
			pemBlock := pem.Block{
				Type:  "CERTIFICATE REQUEST",
				Bytes: csrBytes,
			}
			exit.OnErr(pem.Encode(os.Stdout, &pemBlock))
		},
	}

	subject = setPkixNameFlags(cmd)
	cmd.Flags().StringSliceVar(&dnsNames, "dns-names", nil, "Comma separated list of dns names")
	cmd.Flags().StringSliceVar(&emails, "emails", nil, "Comma separated list of email addresses")
	cmd.Flags().IPSliceVar(&ips, "ip-addrs", nil, "Comma separated list of IP addresses")
	return cmd
}
