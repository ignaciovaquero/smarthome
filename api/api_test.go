package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetLogger(t *testing.T) {
	testCases := []struct {
		name     string
		logger   Logger
		expected *API
	}{
		{
			name: "Testing non setting anything",
			expected: &API{
				Logger:    &DefaultLogger{},
				TableName: DefaultTableName,
			},
			logger: nil,
		},
		{
			name: "Testing setting a default Logger",
			expected: &API{
				Logger:    &DefaultLogger{},
				TableName: DefaultTableName,
			},
			logger: &DefaultLogger{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			actual := NewAPI(SetLogger(tc.logger))
			assert.Equal(tt, tc.expected, actual)
		})
	}
}

func TestSetTableName(t *testing.T) {
	testCases := []struct {
		name      string
		tableName func() *string
		expected  *API
	}{
		{
			name: "Testing setting a custom name",
			tableName: func() *string {
				s := "custom"
				return &s
			},
			expected: &API{
				Logger:    &DefaultLogger{},
				TableName: "custom",
			},
		},
		{
			name: "Testing setting an empty name",
			tableName: func() *string {
				s := ""
				return &s
			},
			expected: &API{
				Logger:    &DefaultLogger{},
				TableName: DefaultTableName,
			},
		},
		{
			name:      "Testing non setting anything",
			tableName: func() *string { return nil },
			expected: &API{
				Logger:    &DefaultLogger{},
				TableName: DefaultTableName,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			var actual *API
			if tc.tableName() != nil {
				actual = NewAPI(SetTableName(*tc.tableName()))
			} else {
				actual = NewAPI()
			}
			assert.Equal(tt, tc.expected, actual)
		})
	}
}
