package main

import (
	"testing"
	"time"

	"snippetbox.brainwhat/internal/assert"
)

func TestHumanDate(t *testing.T) {
	// a slice of anonymous structs
	tests := []struct {
		name     string
		date     time.Time
		expected string
	}{
		{
			name:     "UTC",
			date:     time.Date(2025, 1, 29, 15, 1, 0, 0, time.UTC),
			expected: "29 Jan 2025 at 15:01",
		},
		{
			name:     "Empty",
			date:     time.Time{},
			expected: "",
		},
		{
			name:     "UTC+3",
			date:     time.Date(2025, 01, 29, 15, 1, 0, 0, time.FixedZone("UTC+3", 3*60*60)),
			expected: "29 Jan 2025 at 12:01",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := humanDate(test.date)
			assert.Equal(t, got, test.expected)
		})
	}
}
