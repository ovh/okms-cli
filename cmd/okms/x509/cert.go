package x509

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/ovh/okms-cli/cmd/okms/common"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/spf13/cobra"
)

func createGenerateCertCommand() *cobra.Command {
	var (
		subject  *pkixNameParams
		validity time.Duration
		dnsNames []string
		emails   []string
		ips      []net.IP

		usageServerAuth bool
		usageClientAuth bool
	)

	cmd := &cobra.Command{
		Use:   "cert KEY-ID",
		Short: "Generate a self-signed certificate, signed with the key identified by KEY-ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			keyId := exit.OnErr2(uuid.Parse(args[0]))
			signer := exit.OnErr2(common.Client().NewSigner(cmd.Context(), keyId))

			// Certificate template
			certTemplate := &x509.Certificate{
				// SignatureAlgorithm:
				SerialNumber:          newSerialNumber(),
				IsCA:                  false,
				Subject:               subject.pkixName(),
				SubjectKeyId:          keyId[:],
				NotBefore:             time.Now(),
				NotAfter:              time.Now().Add(validity),
				KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
				BasicConstraintsValid: false,

				DNSNames:       dnsNames,
				EmailAddresses: emails,
				IPAddresses:    ips,
				URIs:           []*url.URL{}, //TODO: Make it a configurable option ?

				// Other fields that can be configured

				// AuthorityKeyId
				// CRLDistributionPoints
				// ExcludedDNSDomains
				// ExcludedEmailAddresses
				// ExcludedIPRanges
				// ExcludedURIDomains
				// ExtKeyUsage
				// ExtraExtensions
				// IssuingCertificateURL
				// MaxPathLen
				// MaxPathLenZero
				// OCSPServer
				// PermittedDNSDomains
				// PermittedDNSDomainsCritical
				// PermittedEmailAddresses
				// PermittedIPRanges
				// PermittedURIDomains
				// PolicyIdentifiers
				// UnknownExtKeyUsage
			}

			if usageClientAuth {
				certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
			}
			if usageServerAuth {
				certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
			}

			cert := exit.OnErr2(x509.CreateCertificate(rand.Reader, certTemplate, certTemplate, signer.Public(), signer))
			pemBlock := pem.Block{
				Type:  "CERTIFICATE",
				Bytes: cert,
			}
			exit.OnErr(pem.Encode(os.Stdout, &pemBlock))
		},
	}

	subject = setPkixNameFlags(cmd)
	cmd.Flags().StringSliceVar(&dnsNames, "dns-names", nil, "Comma separated list of dns names")
	cmd.Flags().StringSliceVar(&emails, "emails", nil, "Comma separated list of email addresses")
	cmd.Flags().IPSliceVar(&ips, "ip-addrs", nil, "Comma separated list of IP addresses")

	cmd.Flags().BoolVar(&usageServerAuth, "server-auth", false, "Enable server auth extended key usage")
	cmd.Flags().BoolVar(&usageServerAuth, "client-auth", false, "Enable client auth extended key usage")

	cmd.Flags().DurationVar(&validity, "validity", 365*24*time.Hour, "Validity duration")

	return cmd
}
