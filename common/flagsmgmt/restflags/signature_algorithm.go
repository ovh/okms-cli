package restflags

import (
	"crypto"
	"errors"
	"hash"

	"github.com/ovh/okms-sdk-go/types"
)

type SignatureAlgorithm string

const (
	ES256 SignatureAlgorithm = "ES256"
	ES384 SignatureAlgorithm = "ES384"
	ES512 SignatureAlgorithm = "ES512"
	RS256 SignatureAlgorithm = "RS256"
	RS384 SignatureAlgorithm = "RS384"
	RS512 SignatureAlgorithm = "RS512"
	PS256 SignatureAlgorithm = "PS256"
	PS384 SignatureAlgorithm = "PS384"
	PS512 SignatureAlgorithm = "PS512"
)

func (e *SignatureAlgorithm) String() string {
	return string(*e)
}

func (e *SignatureAlgorithm) Set(v string) error {
	switch v {
	case "ES256", "ES384", "ES512", "RS256", "RS384", "RS512", "PS256", "PS384", "PS512":
		*e = SignatureAlgorithm(v)
		return nil
	default:
		return errors.New(`must be one of "ES256", "ES384", "ES512", "RS256", "RS384", "RS512", "PS256", "PS384", "PS512"`)
	}
}

func (e *SignatureAlgorithm) Type() string {
	return "ES256|ES384|ES512|RS256|RS384|RS512|PS256|PS384|PS512"
}

func (e *SignatureAlgorithm) HashAlgorithm() crypto.Hash {
	switch *e {
	case ES256, RS256, PS256:
		return crypto.SHA256
	case ES384, RS384, PS384:
		return crypto.SHA384
	case ES512, RS512, PS512:
		return crypto.SHA512
	default:
		panic("Unsupported algorithm:" + *e)
	}
}

func (e *SignatureAlgorithm) Alg() types.DigitalSignatureAlgorithms {
	return types.DigitalSignatureAlgorithms(*e)
}

func (e *SignatureAlgorithm) NewHasher() hash.Hash {
	return e.HashAlgorithm().New()
}
