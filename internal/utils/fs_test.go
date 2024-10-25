package utils

import (
	"os"
	"testing"
)

func TestExpandTilde(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skipf("Cannot get home directory:%s", err.Error())
	}
	tcs := []struct{ in, expected string }{
		{"", ""},
		{"/foo/bar", "/foo/bar"},
		{"~foo", "~foo"},
		{"foo~", "foo~"},
		{"/~foo~/~bar~/~", "/~foo~/~bar~/~"},
		{"~/", home},
		{"~", home},
		{"~/~", home + "/~"},
	}
	for _, tc := range tcs {
		result, err := ExpandTilde(tc.in)
		if err != nil {
			t.Fatalf("ExpandTilde call failed: %s", err.Error())
		}
		if result != tc.expected {
			t.Logf("Input %q gave output %q. Expected %q", tc.in, result, tc.expected)
			t.Fail()
		}
	}
}

func TestExpandTildeError(t *testing.T) {
	for _, env := range []string{"HOME", "home", "USERPROFILE"} {
		_ = os.Setenv(env, "")
	}
	if _, err := ExpandTilde("/ok"); err != nil {
		t.Fatalf("Expected no error, got %q", err.Error())
	}

	if _, err := ExpandTilde("~"); err == nil {
		t.Fatalf("Expected an error, but got nil")
	}
}
