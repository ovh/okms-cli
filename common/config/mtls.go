package config

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/ovh/okms-cli/common/utils/x509utils"
	"github.com/ovh/okms-cli/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	RegisterAuthMethod("", newMtlsFileAuth) // Default to mtls
	RegisterAuthMethod("mtls", newMtlsFileAuth)
}

func newMtlsFileAuth(cmd *cobra.Command, k *koanf.Koanf, envPrefix string) EndpointAuth {
	certFile := GetString(k, "cert", envPrefix+"_CERT", cmd.Flags().Lookup("cert"))
	if certFile == "" {
		exit.OnErr(errors.New("Missing certificate file parameter"))
	}
	keyFile := GetString(k, "key", envPrefix+"_KEY", cmd.Flags().Lookup("key"))
	if keyFile == "" {
		exit.OnErr(errors.New("Missing private key parameter"))
	}

	cert, err := tls.LoadX509KeyPair(exit.OnErr2(utils.ExpandTilde(certFile)), exit.OnErr2(utils.ExpandTilde(keyFile)))
	if err != nil {
		exit.OnErr(fmt.Errorf("Failed to load certificate: %w", err))
	}

	auth := &mtlsFileAuth{certs: []tls.Certificate{cert}}

	okmsIdStr, err := getOkmsId(cert.Leaf)
	if err != nil || okmsIdStr == "*" {
		okmsIdStr = GetString(k, "okmsId", envPrefix+"_OKMSID", cmd.Flags().Lookup("okmsId"))
	}
	auth.okmsId, _ = uuid.Parse(okmsIdStr)

	return auth
}

type mtlsFileAuth struct {
	certs  []tls.Certificate
	okmsId uuid.UUID
}

func (epCfg *mtlsFileAuth) GetToken() *string {
	return nil
}

func (epCfg *mtlsFileAuth) GetOkmsId() uuid.UUID {
	if epCfg.okmsId == uuid.Nil {
		exit.OnErr(errors.New("Invalid OKMS ID"))
	}
	return epCfg.okmsId
}

func (epCfg *mtlsFileAuth) TlsCertificates() []tls.Certificate {
	return epCfg.certs
}

func loadCertPool(caFile string) *x509.CertPool {
	return exit.OnErr2(x509utils.LoadCertPool(caFile))
}

func getOkmsId(cert *x509.Certificate) (string, error) {
	for _, ext := range cert.Extensions { //
		// See https://datatracker.ietf.org/doc/html/rfc5280#section-4.2.1.6
		if !ext.Id.Equal(asn1.ObjectIdentifier{2, 5, 29, 17}) {
			continue
		}
		var seq asn1.RawValue
		_, err := asn1.Unmarshal(ext.Value, &seq)
		if err != nil {
			return "", err
		}
		for rest := seq.Bytes; len(rest) > 0; {
			var val asn1.RawValue
			rest, err = asn1.Unmarshal(rest, &val)
			if err != nil {
				return "", err
			}
			if val.Tag != 0 {
				continue
			}

			var oid asn1.ObjectIdentifier
			rem, err := asn1.Unmarshal(val.Bytes, &oid)
			if err != nil {
				return "", err
			}
			if !oid.Equal(asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 20, 2, 3}) {
				continue
			}
			if _, err = asn1.Unmarshal(rem, &val); err != nil {
				return "", err
			}
			var othername string
			if _, err := asn1.Unmarshal(val.Bytes, &othername); err != nil {
				return "", err
			}
			prefix := "okms.domain:"
			if strings.HasPrefix(othername, prefix) {
				return othername[len(prefix):], nil
			}
		}
	}
	return "", errors.New("No OKMS domain id found")
}
