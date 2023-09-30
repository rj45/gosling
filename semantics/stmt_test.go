package semantics_test

import (
	"strings"
	"testing"

	"github.com/rj45/gosling/types"
)

func TestTypeCheckingStmts(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
		err      string
	}{
		{
			name:     "for condition needs bool",
			src:      "for 1 {}",
			expected: "",
			err:      "for condition must be bool",
		},
		{
			name:     "for ;; condition needs bool",
			src:      "for ;1; {}",
			expected: "",
			err:      "for condition must be bool",
		},
		{
			name:     "assign untyped int converted to int",
			src:      "a := 1",
			expected: "int",
			err:      "",
		},
		{
			name:     "reassign true to false",
			src:      "a := true; a = false",
			expected: "bool",
			err:      "",
		},
		{
			name:     "reassign int to bool",
			src:      "a := 1; a = true",
			expected: "",
			err:      "cannot assign bool to int",
		},
		{
			name:     "variable redefinition",
			src:      "a := 1; a := true",
			expected: "",
			err:      "cannot redefine a",
		},
		{
			name:     "variable used in wrong scope",
			src:      "{a := 1}; a",
			expected: "",
			err:      "undefined name a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, uni, node, errs := parseStmt(t, tt.src)

			if errs != nil {
				for _, err := range errs {
					if tt.err == "" {
						t.Errorf("Expected no error, but got %s", errs)
					} else if !strings.Contains(err.Error(), tt.err) {
						t.Errorf("Expected error to contain %q, but got %q", tt.err, err)
					}
				}
				return
			}

			if errs == nil && tt.err != "" {
				t.Errorf("Expected error %q, but got none", tt.err)
				return
			}

			actual := a.Type(node)

			if actual == types.None && tt.expected != "" {
				t.Errorf("Expected type %q, but got none", tt.expected)
			} else if uni.StringOf(actual) != tt.expected {
				t.Errorf("Expected: %s\nBut got: %s", tt.expected, uni.StringOf(actual))
			}
		})
	}
}
