package util

import (
	"testing"
)

func TestConvertTime(t *testing.T) {
	// Test cases
	testCases := []struct {
		input    string
		format   string
		expected string
	}{
		// Test valid time conversions
		{"2023-09-18 12:34:56", "2006-01-02 15:04:05", "2023-09-18 20:34:56"},
		{"2023-01-01 00:00:00", "2006-01-02 15:04:05", "2023-01-01 08:00:00"},

		// Test invalid time format
		{"2023-09-18 12:34:56", "invalid-format", ""},
	}

	for _, testCase := range testCases {
		result, _, err := ConvertTime(testCase.input, testCase.format)

		if err != nil && testCase.expected != "" {
			t.Errorf("Expected no error for input %s, but got an error: %v", testCase.input, err)
		}

		if result != testCase.expected {
			t.Errorf("For input %s and format %s, expected %s, but got %s", testCase.input, testCase.format, testCase.expected, result)
		}
	}
}
