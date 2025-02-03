package kmipflags

import (
	"errors"
	"strings"

	"github.com/ovh/kmip-go"
	"github.com/ovh/kmip-go/ttlv"
)

type RevocationReason kmip.RevocationReasonCode

const (
	Unspecified          = RevocationReason(kmip.RevocationReasonCodeUnspecified)
	KeyCompromise        = RevocationReason(kmip.RevocationReasonCodeKeyCompromise)
	CACompromise         = RevocationReason(kmip.RevocationReasonCodeCACompromise)
	AffiliationChanged   = RevocationReason(kmip.RevocationReasonCodeAffiliationChanged)
	Superseded           = RevocationReason(kmip.RevocationReasonCodeSuperseded)
	CessationOfOperation = RevocationReason(kmip.RevocationReasonCodeCessationOfOperation)
	PrivilegeWithdrawn   = RevocationReason(kmip.RevocationReasonCodePrivilegeWithdrawn)
)

func (e *RevocationReason) String() string {
	return ttlv.EnumStr(kmip.RevocationReasonCode(*e))
}

func (e *RevocationReason) Set(v string) error {
	switch strings.ToLower(v) {
	case "unspecified":
		*e = Unspecified
	case "keycompromise":
		*e = KeyCompromise
	case "cacompromise":
		*e = CACompromise
	case "affiliationchanged":
		*e = AffiliationChanged
	case "superseded":
		*e = Superseded
	case "cessationofoperation":
		*e = CessationOfOperation
	case "privilegewithdrawn":
		*e = PrivilegeWithdrawn
	default:
		return errors.New(`must be one of "Unspecified", "KeyCompromise", "CACompromise", "AffiliationChanged", "Superseded", "CessationOfOperation", "PrivilegeWithdrawn"`)
	}
	return nil
}

func (e *RevocationReason) Type() string {
	return "Unspecified|KeyCompromise|CACompromise|AffiliationChanged|Superseded|CessationOfOperation|PrivilegeWithdrawn"
}
