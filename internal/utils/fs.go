package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExpandTilde expands the filepath 'in' by replacing the leading tilde ('~') if any by
// the full path to user home directory. It returns an error if the home directory
// cannot be located.
func ExpandTilde(in string) (string, error) {
	if in == "~" || strings.HasPrefix(in, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to expand tilde in %q: %w", in, err)
		}
		return filepath.Join(home, in[1:]), nil
	}
	return in, nil
}
