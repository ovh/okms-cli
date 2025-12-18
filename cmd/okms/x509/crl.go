package x509

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/ovh/okms-cli/cmd/okms/common"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/ovh/okms-cli/common/utils/x509utils"
	"github.com/spf13/cobra"
)

type revocationEntry struct {
	SerialNumber   x509utils.BigInt `json:"serialNumber"`
	RevocationDate time.Time        `json:"revocationDate"`
	ReasonCode     *int             `json:"reasonCode,omitempty"`
}

func CreateGenerateCrlCommand() *cobra.Command {
	var (
		nextUpdate time.Duration
		crlNumber  int64
	)

	cmd := &cobra.Command{
		Use:   "crl CA REVOKE_LIST [KEY-ID]",
		Short: "Generate a CRL with a CA whose key is stored in the KMS",
		Long:  "Generate a Certificate Revocation List with a Certificate Authority whose key is stored in the KMS.\nThe REVOKE_LIST file is a JSON array of entries containing serialNumber (prefer a decimal string; hex must be 0x-prefixed if used), revocationDate and optionally reasonCode. See RFC3339.",
		Args:  cobra.RangeArgs(2, 3),
		Run: func(cmd *cobra.Command, args []string) {
			caData := exit.OnErr2(os.ReadFile(args[0]))
			caDer, _ := pem.Decode(caData)
			ca := exit.OnErr2(x509.ParseCertificate(caDer.Bytes))

			var revocationEntries []revocationEntry
			revoke := exit.OnErr2(os.ReadFile(args[1]))
			exit.OnErr(json.Unmarshal(revoke, &revocationEntries))

			var keyId uuid.UUID
			if len(args) >= 3 {
				keyId = exit.OnErr2(uuid.Parse(args[2]))
			} else {
				if len(ca.SubjectKeyId) == 0 || len(ca.SubjectKeyId) != 16 {
					exit.OnErr(errors.New("Cannot use CA's subject key id, please provide the KEY-ID"))
				}
				keyId = uuid.UUID(ca.SubjectKeyId)
			}

			crl := &x509.RevocationList{
				Number:             big.NewInt(crlNumber),
				SignatureAlgorithm: ca.SignatureAlgorithm,
				Issuer:             ca.Subject,
				ThisUpdate:         time.Now(),
				NextUpdate:         time.Now().Add(nextUpdate),
			}

			for _, entry := range revocationEntries {
				e := x509.RevocationListEntry{
					SerialNumber:   entry.SerialNumber.Int,
					RevocationTime: entry.RevocationDate,
				}

				if entry.ReasonCode != nil {
					e.ReasonCode = *entry.ReasonCode
				}

				crl.RevokedCertificateEntries = append(crl.RevokedCertificateEntries, e)
			}

			signer := exit.OnErr2(common.Client().NewSigner(cmd.Context(), common.GetOkmsId(), keyId))
			certBytes := exit.OnErr2(x509.CreateRevocationList(rand.Reader, crl, ca, signer))

			pemBlock := pem.Block{
				Type:  "X509 CRL",
				Bytes: certBytes,
			}
			exit.OnErr(pem.Encode(os.Stdout, &pemBlock))
		},
	}

	cmd.Flags().DurationVar(&nextUpdate, "nextUpdate", 30*24*time.Hour, "Duration before next update of the CRL, see RFC3339")
	cmd.Flags().Int64Var(&crlNumber, "crlNumber", 1, "CRL Number i.e version, see RFC3339")

	return cmd
}
