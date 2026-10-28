package common

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/ovh/okms-cli/common/config"
	"github.com/ovh/okms-sdk-go/types"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompleteKeyIDs(t *testing.T) {
	const (
		okmsID = "11111111-1111-1111-1111-111111111111"
		keyID  = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/"+okmsID+"/v1/servicekey", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"is_truncated": false,
			"objects_list": []map[string]any{{
				"id":              keyID,
				"name":            "my key",
				"type":            "oct",
				"protectionLevel": "software",
			}},
		})
	}))
	defer server.Close()

	t.Setenv("HOME", t.TempDir())

	cmd := &cobra.Command{Use: "test"}
	config.SetupConfigFlags(cmd)
	SetupRestApiFlags(cmd, nil)
	require.NoError(t, cmd.PersistentFlags().Set("endpoint", server.URL))
	require.NoError(t, cmd.PersistentFlags().Set("auth-method", "token"))
	require.NoError(t, cmd.PersistentFlags().Set("token", "test-token"))
	require.NoError(t, cmd.PersistentFlags().Set("okmsId", okmsID))
	require.NoError(t, cmd.ParseFlags(nil))

	t.Cleanup(func() {
		restClient = nil
		okmsId = uuid.Nil
	})

	completions, directive := CompleteKeyIDs(cmd, nil, "")

	assert.Equal(t, []string{keyID + "\tmy key (oct)"}, completions)
	assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)

	// Key IDs already given on the command line are not suggested again.
	completions, _ = CompleteKeyIDArgs(cmd, []string{keyID}, "")
	assert.Empty(t, completions)

	// Key IDs are only suggested at the expected argument position.
	completions, directive = CompleteKeyIDAt(1)(cmd, nil, "")
	assert.Empty(t, completions)
	assert.Equal(t, cobra.ShellCompDirectiveDefault, directive)

	// Only the keys matching the filter are suggested.
	completions, _ = CompleteKeyIDAtFiltered(0, func(*types.GetServiceKeyResponse) bool { return false })(cmd, nil, "")
	assert.Empty(t, completions)

	completions, _ = CompleteKeyIDAtFiltered(0, func(*types.GetServiceKeyResponse) bool { return true })(cmd, nil, "")
	assert.Equal(t, []string{keyID + "\tmy key (oct)"}, completions)
}
