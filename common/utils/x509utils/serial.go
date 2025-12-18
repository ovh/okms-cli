package x509utils

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
)

type BigInt struct {
	*big.Int
}

func (b BigInt) MarshalJSON() ([]byte, error) {
	if b.Int == nil {
		return []byte("null"), nil
	}
	return json.Marshal(b.String())
}

func (b *BigInt) UnmarshalJSON(data []byte) error {
	if b == nil {
		return fmt.Errorf("nil receiver")
	}
	if string(data) == "null" {
		b.Int = nil
		return nil
	}
	// try string
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		// support decimal and optional 0x/0X-prefixed hexadecimal inputs
		if strings.HasPrefix(strings.ToLower(s), "0x") {
			z, ok := new(big.Int).SetString(s[2:], 16)
			if !ok {
				return fmt.Errorf("invalid hex big.Int string")
			}
			b.Int = z
			return nil
		}
		// decimal only (do not accept plain hex without 0x prefix)
		if z, ok := new(big.Int).SetString(s, 10); ok {
			b.Int = z
			return nil
		}
		return fmt.Errorf("invalid decimal big.Int string")
	}
	// try number (unquoted)
	var n json.Number
	if err := json.Unmarshal(data, &n); err == nil {
		z, ok := new(big.Int).SetString(n.String(), 10)
		if !ok {
			return fmt.Errorf("invalid big.Int number")
		}
		b.Int = z
		return nil
	}
	return fmt.Errorf("invalid big.Int JSON")
}
