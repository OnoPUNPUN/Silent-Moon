package database

import "testing"

func TestNormalizeMigrationURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "postgresql scheme",
			input:    "postgresql://user:pass@host:5432/dbname?sslmode=require",
			expected: "pgx5://user:pass@host:5432/dbname?sslmode=require",
		},
		{
			name:     "postgres scheme",
			input:    "postgres://user:pass@host:5432/dbname?sslmode=require",
			expected: "pgx5://user:pass@host:5432/dbname?sslmode=require",
		},
		{
			name:     "pgx5 scheme already set",
			input:    "pgx5://user:pass@host:5432/dbname?sslmode=require",
			expected: "pgx5://user:pass@host:5432/dbname?sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeMigrationURL(tt.input)
			if got != tt.expected {
				t.Errorf("normalizeMigrationURL(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
