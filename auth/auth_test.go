package auth

import (
	"strings"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name    string
		inputID uint
	}{
		{
			name:    "generate token for ID 1",
			inputID: 1,
		},
		{
			name:    "generate token for ID 10",
			inputID: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.inputID)
			if err != nil {
				t.Fatalf("GenerateToken() error = %v", err)
			}
			if token == "" {
				t.Errorf("GenerateToken() returned empty token")
			}
			// Optionally check if token roughly looks like JWT (three parts separated by dots)
			parts := len(splitToken(token))
			if parts != 3 {
				t.Errorf("GenerateToken() returned invalid JWT format, parts = %d", parts)
			}
		})
	}
}

// Helper to split JWT token by '.'
func splitToken(token string) []string {
	return strings.Split(token, ".")
}
