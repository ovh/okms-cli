package restflags

import (
	"errors"

	"github.com/ovh/okms-sdk-go/types"
)

type CurveType string

const (
	CURVE_P256 CurveType = "P-256"
	CURVE_P384 CurveType = "P-384"
	CURVE_P521 CurveType = "P-521"
)

func (e *CurveType) String() string {
	return string(*e)
}

func (e *CurveType) Set(v string) error {
	switch v {
	case "P-256", "P-384", "P-521":
		*e = CurveType(v)
		return nil
	default:
		return errors.New(`must be one of "P-256", "P-384", "P-521"`)
	}
}

func (e *CurveType) ToRestCurve() types.Curves {
	return types.Curves(*e)
}

func (e *CurveType) Type() string {
	return "P-256|P-384|P-521"
}
