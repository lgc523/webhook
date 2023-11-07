package resp

import (
	"testing"
)

func TestDefaultMsg(t *testing.T) {
	// Test the DefaultMsg function
	expected := "阿里嘎多"
	result := DefaultMsg()

	if result != expected {
		t.Errorf("Expected default message '%s', but got '%s'", expected, result)
	}
}

func TestMsg(t *testing.T) {
	// Test the Msg function for known codes
	testCases := []struct {
		code     int
		expected string
	}{
		{SUCCESS, "OK"},
		{ILLEGAL, "请求不合法"},
		{NotExist, "资源不存在"},
		{DEFAULT, "阿里嘎多"},
	}

	for _, testCase := range testCases {
		result := Msg(testCase.code)

		if result != testCase.expected {
			t.Errorf("For code %d, expected message '%s', but got '%s'", testCase.code, testCase.expected, result)
		}
	}

	// Test the Msg function for an unknown code
	unknownCode := 999
	expectedUnknown := "阿里嘎多"
	resultUnknown := Msg(unknownCode)

	if resultUnknown != expectedUnknown {
		t.Errorf("For unknown code %d, expected default message '%s', but got '%s'", unknownCode, expectedUnknown, resultUnknown)
	}
}
