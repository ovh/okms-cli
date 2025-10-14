package x509

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/ovh/okms-cli/cmd/okms/common"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/spf13/cobra"
)

func createSignCsrCommand() *cobra.Command {
	var (
		validity        time.Duration
		usageServerAuth bool
		usageClientAuth bool

		isCA bool
	)

	cmd := &cobra.Command{
		Use:   "sign CSR CA [KEY-ID]",
		Short: "Sign a certificate request with a CA whose key is stored in the KMS",
		Long: `Sign a certificate request with a CA whose key is stored in the KMS.

The KEY-ID parameter can be left empty if the CA's Subject Key Id matches the key id UUID. Otherwise,
KEY-ID must be the CA's private key UUID`,
		Args: cobra.RangeArgs(2, 3),
		Run: func(cmd *cobra.Command, args []string) {
			csrData := exit.OnErr2(os.ReadFile(args[0]))
			caData := exit.OnErr2(os.ReadFile(args[1]))

			caDer, _ := pem.Decode(caData)
			ca := exit.OnErr2(x509.ParseCertificate(caDer.Bytes))

			csrDer, _ := pem.Decode(csrData)
			csr := exit.OnErr2(x509.ParseCertificateRequest(csrDer.Bytes))
			exit.OnErr(csr.CheckSignature())

			var keyId uuid.UUID
			if len(args) >= 3 {
				keyId = exit.OnErr2(uuid.Parse(args[2]))
			} else {
				if len(ca.SubjectKeyId) == 0 || len(ca.SubjectKeyId) != 16 {
					exit.OnErr(errors.New("Cannot use CA's subject key id, please provide the KEY-ID"))
				}
				keyId = uuid.UUID(ca.SubjectKeyId)
			}

			signer := exit.OnErr2(common.Client().NewSigner(cmd.Context(), common.GetOkmsId(), keyId))

			certTemplate := &x509.Certificate{
				SignatureAlgorithm: ca.SignatureAlgorithm,
				DNSNames:           csr.DNSNames,
				EmailAddresses:     csr.EmailAddresses,
				IPAddresses:        csr.IPAddresses,
				URIs:               csr.URIs,
				ExtraExtensions:    csr.Extensions,

				IsCA:                  isCA,
				BasicConstraintsValid: !isCA,

				SerialNumber:   newSerialNumber(),
				Issuer:         ca.Subject,
				Subject:        csr.Subject,
				AuthorityKeyId: keyId[:],
				NotBefore:      time.Now(),
				ExtKeyUsage:    []x509.ExtKeyUsage{},
			}
			if isCA {
				certTemplate.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature
			} else {
				certTemplate.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement
				if usageClientAuth {
					certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
				}
				if usageServerAuth {
					certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
				}
			}

			if csrKeyId, found := extractCsrSubjectKeyId(csr); found {
				certTemplate.SubjectKeyId = csrKeyId[:]
			}

			certBytes := exit.OnErr2(x509.CreateCertificate(rand.Reader, certTemplate, ca, csr.PublicKey, signer))

			pemBlock := pem.Block{
				Type:  "CERTIFICATE",
				Bytes: certBytes,
			}
			exit.OnErr(pem.Encode(os.Stdout, &pemBlock))
		},
	}
	cmd.Flags().DurationVar(&validity, "validity", 365*24*time.Hour, "Validity duration")
	cmd.Flags().BoolVar(&usageServerAuth, "server-auth", false, "Enable server auth extended key usage")
	cmd.Flags().BoolVar(&usageServerAuth, "client-auth", false, "Enable client auth extended key usage")

	cmd.Flags().BoolVar(&isCA, "new-ca", false, "Sign as a CA certificate")

	cmd.MarkFlagsMutuallyExclusive("new-ca", "server-auth")
	cmd.MarkFlagsMutuallyExclusive("new-ca", "client-auth")

	return cmd
}

func extractCsrSubjectKeyId(csr *x509.CertificateRequest) (uuid.UUID, bool) {
	for _, ext := range csr.Extensions {
		if !ext.Id.Equal(OID_CE_SUBJECT_KEY_IDENTIFIER) {
			continue
		}
		var uid []byte
		if _, err := asn1.Unmarshal(ext.Value, &uid); err != nil {
			continue
		}
		validated, err := uuid.FromBytes(uid)
		if err != nil {
			continue
		}
		return validated, true
	}
	return uuid.Nil, false
}
