package restflags

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ovh/okms-sdk-go/types"
)

type RevocationReason types.RevocationReasons

const (
	AffiliationChanged   RevocationReason = RevocationReason(types.AffiliationChanged)
	CaCompromise         RevocationReason = RevocationReason(types.CaCompromise)
	CessationOfOperation RevocationReason = RevocationReason(types.CessationOfOperation)
	KeyCompromise        RevocationReason = RevocationReason(types.KeyCompromise)
	PrivilegeWithdrawn   RevocationReason = RevocationReason(types.PrivilegeWithdrawn)
	Superseded           RevocationReason = RevocationReason(types.Superseded)
	Unspecified          RevocationReason = RevocationReason(types.Unspecified)
)

var revocationReasonValues = []string{
	string(AffiliationChanged),
	string(CaCompromise),
	string(CessationOfOperation),
	string(KeyCompromise),
	string(PrivilegeWithdrawn),
	string(Superseded),
	string(Unspecified),
}

func (e *RevocationReason) String() string {
	return string(*e)
}

func (e *RevocationReason) Set(v string) error {
	if slices.ContainsFunc(revocationReasonValues, func(e string) bool { return e == v }) {
		*e = RevocationReason(v)
		return nil
	}
	return fmt.Errorf("must be one of %s", strings.Join(revocationReasonValues, ", "))
}

func (e *RevocationReason) Type() string {
	return strings.Join(revocationReasonValues, "|")
}

func (e RevocationReason) RestModel() types.RevocationReasons {
	return types.RevocationReasons(e)
}
