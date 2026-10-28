package common

import (
	"context"
	"fmt"
	"strings"

	"github.com/ovh/okms-sdk-go/types"
	"github.com/spf13/cobra"
)

// completionDescReplacer sanitizes the values used as completion descriptions,
// as the shell completion protocol is line based.
var completionDescReplacer = strings.NewReplacer("\t", " ", "\n", " ", "\r", " ")

// KeyFilter restricts the keys suggested by the key ID completion functions.
// A nil KeyFilter suggests all the keys.
type KeyFilter func(*types.GetServiceKeyResponse) bool

// CompleteKeyIDs is a completion function suggesting the IDs of the domain
// keys. It can be used as a ValidArgsFunction, or registered as a flag
// completion function with RegisterFlagCompletionFunc.
func CompleteKeyIDs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return completeKeyIDs(cmd, nil, nil)
}

// CompleteKeyIDsFiltered is like CompleteKeyIDs, but only suggests the keys
// accepted by filter.
func CompleteKeyIDsFiltered(filter KeyFilter) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completeKeyIDs(cmd, nil, filter)
	}
}

// CompleteKeyIDArgs is like CompleteKeyIDs, but it does not suggest the key IDs
// already given in args. It is meant for commands accepting several KEY-ID
// arguments, such as `keys delete`.
func CompleteKeyIDArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return completeKeyIDs(cmd, args, nil)
}

// CompleteKeyIDAt returns a completion function suggesting key IDs only when
// completing the argument at the given position (0 based). The completion of
// the other arguments falls back to the default behavior (file completion).
func CompleteKeyIDAt(position int) cobra.CompletionFunc {
	return CompleteKeyIDAtFiltered(position, nil)
}

// CompleteKeyIDAtFiltered is like CompleteKeyIDAt, but only suggests the keys
// accepted by filter. It is meant for commands that only accept a subset of the
// keys, such as `keys export`, which exports public key material by default but
// requires an extractable key to export key material in wrapped form.
func CompleteKeyIDAtFiltered(position int, filter KeyFilter) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != position {
			return nil, cobra.ShellCompDirectiveDefault
		}
		return completeKeyIDs(cmd, nil, filter)
	}
}

func completeKeyIDs(cmd *cobra.Command, exclude []string, filter KeyFilter) ([]string, cobra.ShellCompDirective) {
	InitClient(cmd)
	if restClient == nil {
		// No usable configuration: do not fall back to file completion.
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	skipped := make(map[string]struct{}, len(exclude))
	for _, id := range exclude {
		skipped[id] = struct{}{}
	}

	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// List all keys, including deactivated ones, as they are the ones the user
	// may want to activate or delete.
	state := types.KeyStatesAll
	completions := make([]string, 0)
	for key, err := range restClient.ListAllServiceKeys(GetOkmsId(), nil, &state).Iter(ctx) {
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		id := key.Id.String()
		if _, ok := skipped[id]; ok {
			continue
		}
		if filter != nil && !filter(key) {
			continue
		}
		completions = append(completions, cobra.CompletionWithDesc(id, keyCompletionDescription(key)))
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

func keyCompletionDescription(key *types.GetServiceKeyResponse) string {
	desc := key.Name
	if key.Type != "" {
		desc = fmt.Sprintf("%s (%s)", key.Name, key.Type)
	}
	return completionDescReplacer.Replace(desc)
}
