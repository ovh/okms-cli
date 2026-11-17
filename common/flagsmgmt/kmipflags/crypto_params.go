package kmipflags

import (
	"errors"
	"strings"

	"github.com/ovh/kmip-go"
	"github.com/ovh/kmip-go/ttlv"
)

type BlockCipherMode kmip.BlockCipherMode

const (
	GCM BlockCipherMode = BlockCipherMode(kmip.BlockCipherModeGCM)
	CBC BlockCipherMode = BlockCipherMode(kmip.BlockCipherModeCBC)
	ECB BlockCipherMode = BlockCipherMode(kmip.BlockCipherModeECB)
)

func (e *BlockCipherMode) String() string {
	if *e == 0 {
		return ""
	}
	return ttlv.EnumStr(kmip.BlockCipherMode(*e))
}

func (e *BlockCipherMode) Set(v string) error {
	switch strings.ToLower(v) {
	case "gcm":
		*e = GCM
	case "cbc":
		*e = CBC
	case "ecb":
		*e = ECB
	default:
		return errors.New(`must be one of "GCM", "CBC", "ECB"`)
	}
	return nil
}

func (e *BlockCipherMode) Type() string {
	return "GCM|CBC|ECB"
}

type PaddingMethod kmip.PaddingMethod

const (
	PaddingNone     PaddingMethod = PaddingMethod(kmip.PaddingMethodNone)
	PaddingPKCS5    PaddingMethod = PaddingMethod(kmip.PaddingMethodPKCS5)
	PaddingPKCS1v15 PaddingMethod = PaddingMethod(kmip.PaddingMethodPKCS1V1_5)
	PaddingOAEP     PaddingMethod = PaddingMethod(kmip.PaddingMethodOAEP)
)

func (e *PaddingMethod) String() string {
	if *e == 0 {
		return ""
	}
	return ttlv.EnumStr(kmip.PaddingMethod(*e))
}

func (e *PaddingMethod) Set(v string) error {
	switch strings.ToLower(v) {
	case "none":
		*e = PaddingNone
	case "pkcs5", "pkcs7":
		*e = PaddingPKCS5
	case "pkcs1v15", "pkcs1v1_5", "pkcs1":
		*e = PaddingPKCS1v15
	case "oaep":
		*e = PaddingOAEP
	default:
		return errors.New(`must be one of "None", "PKCS5", "PKCS1v15", "OAEP"`)
	}
	return nil
}

func (e *PaddingMethod) Type() string {
	return "None|PKCS5|PKCS1v15|OAEP"
}

type HashingAlgorithm kmip.HashingAlgorithm

const (
	SHA1   HashingAlgorithm = HashingAlgorithm(kmip.HashingAlgorithmSHA_1)
	SHA256 HashingAlgorithm = HashingAlgorithm(kmip.HashingAlgorithmSHA_256)
	SHA384 HashingAlgorithm = HashingAlgorithm(kmip.HashingAlgorithmSHA_384)
	SHA512 HashingAlgorithm = HashingAlgorithm(kmip.HashingAlgorithmSHA_512)
)

func (e *HashingAlgorithm) String() string {
	if *e == 0 {
		return ""
	}
	return ttlv.EnumStr(kmip.HashingAlgorithm(*e))
}

func (e *HashingAlgorithm) Set(v string) error {
	switch strings.ReplaceAll(strings.ToLower(v), "-", "") {
	case "sha1":
		*e = SHA1
	case "sha256":
		*e = SHA256
	case "sha384":
		*e = SHA384
	case "sha512":
		*e = SHA512
	default:
		return errors.New(`must be one of "SHA1", "SHA256", "SHA384", "SHA512"`)
	}
	return nil
}

func (e *HashingAlgorithm) Type() string {
	return "SHA1|SHA256|SHA384|SHA512"
}

// SignatureAlgorithm uses the same names as the REST API signature algorithms.
type SignatureAlgorithm string

const (
	RS256 SignatureAlgorithm = "RS256"
	RS384 SignatureAlgorithm = "RS384"
	RS512 SignatureAlgorithm = "RS512"
	PS256 SignatureAlgorithm = "PS256"
	PS384 SignatureAlgorithm = "PS384"
	PS512 SignatureAlgorithm = "PS512"
	ES256 SignatureAlgorithm = "ES256"
	ES384 SignatureAlgorithm = "ES384"
	ES512 SignatureAlgorithm = "ES512"
)

func (e *SignatureAlgorithm) String() string {
	return string(*e)
}

func (e *SignatureAlgorithm) Set(v string) error {
	switch alg := SignatureAlgorithm(strings.ToUpper(v)); alg {
	case RS256, RS384, RS512, PS256, PS384, PS512, ES256, ES384, ES512:
		*e = alg
		return nil
	default:
		return errors.New(`must be one of "RS256", "RS384", "RS512", "PS256", "PS384", "PS512", "ES256", "ES384", "ES512"`)
	}
}

func (e *SignatureAlgorithm) Type() string {
	return "RS256|RS384|RS512|PS256|PS384|PS512|ES256|ES384|ES512"
}

// CryptographicParameters returns the KMIP cryptographic parameters selecting the signature algorithm:
// RSxxx is RSA PKCS#1 v1.5, PSxxx is RSASSA-PSS and ESxxx is ECDSA, with the SHA-2 hash xxx.
func (e *SignatureAlgorithm) CryptographicParameters() kmip.CryptographicParameters {
	var hash kmip.HashingAlgorithm
	switch *e {
	case RS256, PS256, ES256:
		hash = kmip.HashingAlgorithmSHA_256
	case RS384, PS384, ES384:
		hash = kmip.HashingAlgorithmSHA_384
	case RS512, PS512, ES512:
		hash = kmip.HashingAlgorithmSHA_512
	}

	var alg kmip.DigitalSignatureAlgorithm
	switch *e {
	case RS256:
		alg = kmip.DigitalSignatureAlgorithmSHA_256WithRSAEncryption
	case RS384:
		alg = kmip.DigitalSignatureAlgorithmSHA_384WithRSAEncryption
	case RS512:
		alg = kmip.DigitalSignatureAlgorithmSHA_512WithRSAEncryption
	case PS256, PS384, PS512:
		alg = kmip.DigitalSignatureAlgorithmRSASSA_PSS
	case ES256:
		alg = kmip.DigitalSignatureAlgorithmECDSAWithSHA256
	case ES384:
		alg = kmip.DigitalSignatureAlgorithmECDSAWithSHA384
	case ES512:
		alg = kmip.DigitalSignatureAlgorithmECDSAWithSHA512
	}

	return kmip.CryptographicParameters{
		DigitalSignatureAlgorithm: alg,
		HashingAlgorithm:          hash,
	}
}

type UsageLimitsUnit kmip.UsageLimitsUnit

const (
	UsageLimitsObject UsageLimitsUnit = UsageLimitsUnit(kmip.UsageLimitsUnitObject)
	UsageLimitsByte   UsageLimitsUnit = UsageLimitsUnit(kmip.UsageLimitsUnitByte)
)

func (e *UsageLimitsUnit) String() string {
	if *e == 0 {
		return ""
	}
	return ttlv.EnumStr(kmip.UsageLimitsUnit(*e))
}

func (e *UsageLimitsUnit) Set(v string) error {
	switch strings.ToLower(v) {
	case "object":
		*e = UsageLimitsObject
	case "byte":
		*e = UsageLimitsByte
	default:
		return errors.New(`must be one of "Object", "Byte"`)
	}
	return nil
}

func (e *UsageLimitsUnit) Type() string {
	return "Object|Byte"
}
