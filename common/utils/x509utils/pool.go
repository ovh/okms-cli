package x509utils

import (
	"crypto/x509"
	"fmt"
	"os"
)

func LoadCertPool(caFile string) (*x509.CertPool, error) {
	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("Could not load system certificates pool: %w", err)
	}

	if caFile != "" {
		caBundle, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("Could not load CA file %q: %w", caFile, err)
		}
		pool.AppendCertsFromPEM(caBundle)
	}
	return pool, nil
}
