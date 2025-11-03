package main

import (
	"reflect"
	"testing"

	"github.com/charmbracelet/vhs/parser"
)

func TestCommand(t *testing.T) {
	const numberOfCommands = 29
	if len(parser.CommandTypes) != numberOfCommands {
		t.Errorf("Expected %d commands, got %d", numberOfCommands, len(parser.CommandTypes))
	}

	const numberOfCommandFuncs = 29
	if len(CommandFuncs) != numberOfCommandFuncs {
		t.Errorf("Expected %d commands, got %d", numberOfCommandFuncs, len(CommandFuncs))
	}
}

func TestExecuteSetTheme(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		theme, err := getTheme("  ")
		requireNoErr(t, err)
		requireDefaultTheme(t, theme)
	})
	t.Run("named", func(t *testing.T) {
		theme, err := getTheme("Andromeda")
		requireNoErr(t, err)
		requireNotDefaultTheme(t, theme)
	})
	t.Run("json", func(t *testing.T) {
		theme, err := getTheme(`{"background": "#29283b"}`)
		requireNoErr(t, err)
		requireNotDefaultTheme(t, theme)
		if "#29283b" != theme.Background {
			t.Errorf("wrong background, expected %q, got %q", "#29283b", theme.Background)
		}
	})
	t.Run("suggestion", func(t *testing.T) {
		theme, err := getTheme("cattppuccin latt")
		requireEqualErr(t, err, "invalid `Set Theme \"cattppuccin latt\"`: did you mean \"Catppuccin Latte\"")
		requireDefaultTheme(t, theme)
	})
	t.Run("invalid json", func(t *testing.T) {
		theme, err := getTheme(`{"background`)
		requireErr(t, err)
		requireDefaultTheme(t, theme)
	})
	t.Run("unknown theme", func(t *testing.T) {
		theme, err := getTheme("foobar")
		requireErr(t, err)
		requireDefaultTheme(t, theme)
	})
}

func requireErr(tb testing.TB, err error) {
	tb.Helper()
	if err == nil {
		tb.Fatalf("expected an error, got nil")
	}
}

func requireEqualErr(tb testing.TB, err1 error, err2 string) {
	tb.Helper()
	if err1 == nil {
		tb.Fatalf("expected an error, got nil")
	}
	if err1.Error() != err2 {
		tb.Fatalf("errors do not match: %q != %q", err1.Error(), err2)
	}
}

func requireNoErr(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatalf("expected no error, got: %v", err)
	}
}

func requireDefaultTheme(tb testing.TB, theme Theme) {
	tb.Helper()
	if !reflect.DeepEqual(DefaultTheme, theme) {
		tb.Fatalf("expected theme to be the default theme, got something else: %+v", theme)
	}
}

func requireNotDefaultTheme(tb testing.TB, theme Theme) {
	tb.Helper()
	if reflect.DeepEqual(DefaultTheme, theme) {
		tb.Fatalf("expected theme to be different from the default theme, got the default instead")
	}
}

func TestInterpolateEnvVars(t *testing.T) {
	t.Setenv("NAME", "Emmanuel Goldstein")
	t.Setenv("LANG", "golang")
	t.Setenv("USER", "egoldstein")

	tests := map[string]struct {
		input    string
		expected string
	}{
		"single": {
			input:    "Hello ${{NAME}}",
			expected: "Hello Emmanuel Goldstein",
		},
		"multiple": {
			input:    "${{LANG}}-${{NAME}}-${{LANG}}",
			expected: "golang-Emmanuel Goldstein-golang",
		},
		"missing": {
			input:    "Hi ${{MISSING}}",
			expected: "Hi ",
		},
		"content-after": {
			input:    "Hello, ${{USER}}!!!",
			expected: "Hello, egoldstein!!!",
		},
	}

	for name, tt := range tests {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			if got := interpolateEnvVars(tt.input); got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
