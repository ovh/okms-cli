package x509

import (
	"crypto/rand"
	"crypto/x509/pkix"
	"encoding/asn1"
	"math"
	"math/big"

	"github.com/ovh/okms-cli/cmd/okms/common"
	"github.com/spf13/cobra"
)

var (
	// SubjectKeyIdentifier as defined in [RFC 5280 Section 4.2.1.2].
	//
	//   SubjectKeyIdentifier ::= KeyIdentifier
	//
	// [RFC 5280 Section 4.2.1.2]: https://datatracker.ietf.org/doc/html/rfc5280#section-4.2.1.2
	OID_CE_SUBJECT_KEY_IDENTIFIER = asn1.ObjectIdentifier{2, 5, 29, 14}
)

func CreateX509Command(cust common.CustomizeFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "x509",
		Short: "Generate, and sign x509 certificates",
	}
	common.SetupRestApiFlags(cmd, cust)
	cmd.AddCommand(
		newCreateCommand(),
		createSignCsrCommand(),
	)

	return cmd
}

type pkixNameParams struct {
	commonName string
	orgs       []string
	country    []string
	orgUnit    []string
}

func setPkixNameFlags(cmd *cobra.Command) *pkixNameParams {
	params := new(pkixNameParams)
	cmd.Flags().StringVar(&params.commonName, "cn", "", "Common Name")
	cmd.Flags().StringSliceVar(&params.orgUnit, "ou", nil, "Comma separated Organizational Units")
	cmd.Flags().StringSliceVar(&params.orgs, "org", nil, "Comma separated Organizations")
	cmd.Flags().StringSliceVar(&params.country, "country", nil, "Comma separated Countries")
	return params
}

func (params *pkixNameParams) pkixName() pkix.Name {
	return pkix.Name{
		CommonName:         params.commonName,
		Organization:       params.orgs,
		Country:            params.country,
		OrganizationalUnit: params.orgUnit,
	}
}

// newSerialNumber creates a new random big integer or panic if it fails to do it.
func newSerialNumber() *big.Int {
	serial, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		panic(err)
	}
	return serial
}
