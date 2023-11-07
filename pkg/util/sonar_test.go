package util

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

// MockTransport is a custom http.RoundTripper that allows us to mock HTTP responses.
type MockTransport struct {
	Response *http.Response
	Err      error
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}

func TestSonarResult(t *testing.T) {
	// Mock the HTTP response for testing.
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(
			bytes.NewBufferString(`{
				"measures": [
					{"metric": "bugs", "value": "1"},
					{"metric": "vulnerabilities", "value": "2"},
					{"metric": "code_smells", "value": "3"},
					{"metric": "duplicated_lines_density", "value": "4"},
					{"metric": "coverage", "value": "80.0"},
					{"metric": "ncloc", "value": "1000"}
				]
			}`)),
	}

	// Create a custom HTTP client with the mock transport.
	httpClient := &http.Client{
		Transport: &MockTransport{Response: mockResponse, Err: nil},
	}

	// Call the SonarResult function with a project key.
	projectKey := "your_project_key" // Replace with your actual project key.
	metricValueMap, err := SonarResult(httpClient, projectKey)

	// Check for errors.
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check the values in the metricValueMap.
	expectedMetrics := map[string]string{
		"bugs":                     "1",
		"vulnerabilities":          "2",
		"code_smells":              "3",
		"duplicated_lines_density": "4",
		"coverage":                 "80.0",
		"ncloc":                    "1000",
	}

	for metric, expectedValue := range expectedMetrics {
		if value, ok := metricValueMap[metric]; ok {
			if value != expectedValue {
				t.Errorf("Expected metric '%s' with value '%s', got '%s'", metric, expectedValue, value)
			}
		} else {
			t.Errorf("Expected metric '%s' not found in the result", metric)
		}
	}
}
