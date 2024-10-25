package config

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"

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
	return &mtlsFileAuth{
		CertFile: exit.OnErr2(utils.ExpandTilde(certFile)),
		KeyFile:  exit.OnErr2(utils.ExpandTilde(keyFile)),
	}
}

type mtlsFileAuth struct {
	CertFile string
	KeyFile  string
}

func (epCfg *mtlsFileAuth) TlsCertificates() []tls.Certificate {
	cert, err := tls.LoadX509KeyPair(epCfg.CertFile, epCfg.KeyFile)
	if err != nil {
		exit.OnErr(fmt.Errorf("Failed to load certificates: %w", err))
	}
	return []tls.Certificate{cert}
}

func loadCertPool(caFile string) *x509.CertPool {
	return exit.OnErr2(x509utils.LoadCertPool(caFile))
}
