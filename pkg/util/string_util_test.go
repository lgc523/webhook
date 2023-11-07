package util

import "testing"

func TestFormatStrF2Int(t *testing.T) {
	// Test cases
	testCases := []struct {
		input    string
		expected int
	}{
		{"12.34", 1234},
		{"0.01", 1},
		{"-5.67", -567},
		{"abc", -1}, // Invalid input
	}

	for _, testCase := range testCases {
		result, err := FormatStrF2Int(testCase.input)

		if err != nil && testCase.expected != -1 {
			t.Errorf("Expected no error for input %s, but got an error: %v", testCase.input, err)
		}

		if result != testCase.expected {
			t.Errorf("For input %s, expected %d, but got %d", testCase.input, testCase.expected, result)
		}
	}
}
