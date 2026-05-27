package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabase_NewConnection(t *testing.T) {
	tests := []struct {
		name        string
		connString  string
		expectError bool
	}{
		{
			name:        "invalid connection string",
			connString:  "invalid://connection",
			expectError: true,
		},
		{
			name:        "empty connection string",
			connString:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := NewNoopLogger()
			_, err := NewConnection(context.Background(), tt.connString, log)

			if tt.expectError {
				assert.Error(t, err)
				assert.True(t, len(err.Error()) > 0, "Expected error message to not be empty")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDatabase_GetPool(t *testing.T) {
	db := &Database{
		Pool:   nil,
		logger: NewNoopLogger(),
	}

	result := db.GetPool()
	assert.Nil(t, result)
}

func TestTruncateSQL(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		expected string
	}{
		{
			name:     "short SQL",
			sql:      "SELECT * FROM users",
			expected: "SELECT * FROM users",
		},
		{
			name:     "long SQL gets truncated",
			sql:      "SELECT * FROM users WHERE name = 'very long name that exceeds the limit and should be truncated for safety in logs'",
			expected: "SELECT * FROM users WHERE name = 'very long name that exceeds the limit and should be truncated for ...",
		},
		{
			name:     "exact limit SQL",
			sql:      string(make([]byte, sqlPreviewMaxLen)),
			expected: string(make([]byte, sqlPreviewMaxLen)),
		},
		{
			name:     "empty SQL",
			sql:      "",
			expected: "",
		},
		{
			name:     "single character SQL",
			sql:      "A",
			expected: "A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateSQL(tt.sql)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDatabase_Close(t *testing.T) {
	db := &Database{
		Pool:   nil,
		logger: NewNoopLogger(),
	}

	assert.NotPanics(t, func() {
		db.Close()
	})
}

func TestDatabase_Shutdown(t *testing.T) {
	db := &Database{
		Pool:   nil,
		logger: NewNoopLogger(),
	}

	assert.NotPanics(t, func() {
		err := db.Shutdown(context.TODO())
		assert.NoError(t, err)
	})
}
