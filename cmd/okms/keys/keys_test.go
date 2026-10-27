package keys

import (
	"testing"

	"github.com/ovh/okms-sdk-go/types"
	"github.com/stretchr/testify/assert"
)

func TestIsExtractableKey(t *testing.T) {
	tests := []struct {
		name       string
		attributes *map[string]any
		want       bool
	}{
		{name: "extractable", attributes: &map[string]any{"extractable": true}, want: true},
		{name: "not extractable", attributes: &map[string]any{"extractable": false}, want: false},
		{name: "missing attribute", attributes: &map[string]any{}, want: false},
		{name: "no attributes", attributes: nil, want: false},
		{name: "invalid attribute type", attributes: &map[string]any{"extractable": "true"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := &types.GetServiceKeyResponse{Attributes: tt.attributes}
			assert.Equal(t, tt.want, isExtractableKey(key))
		})
	}
}
