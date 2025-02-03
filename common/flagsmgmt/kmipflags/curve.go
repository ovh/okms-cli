package kmipflags

import (
	"errors"
	"strings"

	"github.com/ovh/kmip-go"
)

type EcCurve kmip.RecommendedCurve

const (
	P256 EcCurve = EcCurve(kmip.P_256)
	P384 EcCurve = EcCurve(kmip.P_384)
	P521 EcCurve = EcCurve(kmip.P_521)
)

func (e *EcCurve) String() string {
	switch *e {
	case 0:
		return ""
	case P256:
		return "P-256"
	case P384:
		return "P-384"
	case P521:
		return "P-521"
	}
	panic("unreachable")
}

func (e *EcCurve) Set(v string) error {
	switch strings.ToLower(v) {
	case "p-256":
		*e = P256
	case "p-384":
		*e = P384
	case "p-521":
		*e = P521
	default:
		return errors.New(`must be one of "P-256", "P-384", "P-521"`)
	}
	return nil
}

func (e *EcCurve) Type() string {
	return "P-256|P-384|P-521"
}
