package testutil

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/yourusername/flintroute/test/functional/pkg/client"
)

// AssertPeerEqual asserts that two peers are equal
func AssertPeerEqual(t *testing.T, expected, actual *client.Peer) {
	t.Helper()

	if expected.Name != actual.Name {
		t.Errorf("Peer name mismatch: expected %s, got %s", expected.Name, actual.Name)
	}
	if expected.IPAddress != actual.IPAddress {
		t.Errorf("Peer IP address mismatch: expected %s, got %s", expected.IPAddress, actual.IPAddress)
	}
	if expected.ASN != actual.ASN {
		t.Errorf("Peer ASN mismatch: expected %d, got %d", expected.ASN, actual.ASN)
	}
	if expected.RemoteASN != actual.RemoteASN {
		t.Errorf("Peer RemoteASN mismatch: expected %d, got %d", expected.RemoteASN, actual.RemoteASN)
	}
	if expected.Enabled != actual.Enabled {
		t.Errorf("Peer Enabled mismatch: expected %v, got %v", expected.Enabled, actual.Enabled)
	}
}

// AssertSessionState asserts that a session has the expected state
func AssertSessionState(t *testing.T, expected string, actual *client.Session) {
	t.Helper()

	if actual.State != expected {
		t.Errorf("Session state mismatch: expected %s, got %s", expected, actual.State)
	}
}

// AssertAlertExists asserts that an alert with the given message exists
func AssertAlertExists(t *testing.T, alerts []*client.Alert, message string) {
	t.Helper()

	for _, alert := range alerts {
		if strings.Contains(alert.Message, message) {
			return
		}
	}

	t.Errorf("Alert with message containing '%s' not found", message)
}

// AssertHTTPStatus asserts that the HTTP status code matches expected
func AssertHTTPStatus(t *testing.T, expected, actual int, body string) {
	t.Helper()

	if expected != actual {
		t.Errorf("HTTP status mismatch: expected %d, got %d. Body: %s", expected, actual, body)
	}
}

// AssertNoError asserts that there is no error
func AssertNoError(t *testing.T, err error, context string) {
	t.Helper()

	if err != nil {
		t.Fatalf("%s: unexpected error: %v", context, err)
	}
}

// AssertError asserts that there is an error
func AssertError(t *testing.T, err error, context string) {
	t.Helper()

	if err == nil {
		t.Fatalf("%s: expected error but got none", context)
	}
}

// AssertErrorContains asserts that the error message contains the expected substring
func AssertErrorContains(t *testing.T, err error, expected string, context string) {
	t.Helper()

	if err == nil {
		t.Fatalf("%s: expected error but got none", context)
	}

	if !strings.Contains(err.Error(), expected) {
		t.Errorf("%s: error message '%s' does not contain '%s'", context, err.Error(), expected)
	}
}

// AssertEqual asserts that two values are equal
func AssertEqual(t *testing.T, expected, actual interface{}, context string) {
	t.Helper()

	if expected != actual {
		t.Errorf("%s: expected %v, got %v", context, expected, actual)
	}
}

// AssertNotEqual asserts that two values are not equal
func AssertNotEqual(t *testing.T, expected, actual interface{}, context string) {
	t.Helper()

	if expected == actual {
		t.Errorf("%s: expected values to be different, but both are %v", context, expected)
	}
}

// AssertTrue asserts that a condition is true
func AssertTrue(t *testing.T, condition bool, message string) {
	t.Helper()

	if !condition {
		t.Errorf("Assertion failed: %s", message)
	}
}

// AssertFalse asserts that a condition is false
func AssertFalse(t *testing.T, condition bool, message string) {
	t.Helper()

	if condition {
		t.Errorf("Assertion failed: %s", message)
	}
}

// AssertNil asserts that a value is nil
func AssertNil(t *testing.T, value interface{}, context string) {
	t.Helper()

	if value != nil {
		t.Errorf("%s: expected nil, got %v", context, value)
	}
}

// AssertNotNil asserts that a value is not nil
func AssertNotNil(t *testing.T, value interface{}, context string) {
	t.Helper()

	if value == nil {
		t.Errorf("%s: expected non-nil value", context)
	}
}

// AssertGreaterThan asserts that actual is greater than expected
func AssertGreaterThan(t *testing.T, expected, actual int, context string) {
	t.Helper()

	if actual <= expected {
		t.Errorf("%s: expected %d to be greater than %d", context, actual, expected)
	}
}

// AssertLessThan asserts that actual is less than expected
func AssertLessThan(t *testing.T, expected, actual int, context string) {
	t.Helper()

	if actual >= expected {
		t.Errorf("%s: expected %d to be less than %d", context, actual, expected)
	}
}

// AssertContains asserts that a string contains a substring
func AssertContains(t *testing.T, haystack, needle string, context string) {
	t.Helper()

	if !strings.Contains(haystack, needle) {
		t.Errorf("%s: '%s' does not contain '%s'", context, haystack, needle)
	}
}

// AssertNotContains asserts that a string does not contain a substring
func AssertNotContains(t *testing.T, haystack, needle string, context string) {
	t.Helper()

	if strings.Contains(haystack, needle) {
		t.Errorf("%s: '%s' should not contain '%s'", context, haystack, needle)
	}
}

// AssertJSONEqual asserts that two JSON strings are equal
func AssertJSONEqual(t *testing.T, expected, actual string, context string) {
	t.Helper()

	var expectedJSON, actualJSON interface{}

	if err := json.Unmarshal([]byte(expected), &expectedJSON); err != nil {
		t.Fatalf("%s: failed to parse expected JSON: %v", context, err)
	}

	if err := json.Unmarshal([]byte(actual), &actualJSON); err != nil {
		t.Fatalf("%s: failed to parse actual JSON: %v", context, err)
	}

	expectedStr := fmt.Sprintf("%v", expectedJSON)
	actualStr := fmt.Sprintf("%v", actualJSON)

	if expectedStr != actualStr {
		t.Errorf("%s: JSON mismatch\nExpected: %s\nActual: %s", context, expectedStr, actualStr)
	}
}

// AssertSliceLength asserts that a slice has the expected length
func AssertSliceLength(t *testing.T, expected int, slice interface{}, context string) {
	t.Helper()

	var length int
	switch v := slice.(type) {
	case []*client.Peer:
		length = len(v)
	case []*client.Session:
		length = len(v)
	case []*client.Alert:
		length = len(v)
	case []*client.ConfigVersion:
		length = len(v)
	default:
		t.Fatalf("%s: unsupported slice type", context)
	}

	if length != expected {
		t.Errorf("%s: expected slice length %d, got %d", context, expected, length)
	}
}

// AssertPeerExists asserts that a peer with the given IP exists in the list
func AssertPeerExists(t *testing.T, peers []*client.Peer, ipAddress string) {
	t.Helper()

	for _, peer := range peers {
		if peer.IPAddress == ipAddress {
			return
		}
	}

	t.Errorf("Peer with IP address %s not found", ipAddress)
}

// AssertPeerNotExists asserts that a peer with the given IP does not exist in the list
func AssertPeerNotExists(t *testing.T, peers []*client.Peer, ipAddress string) {
	t.Helper()

	for _, peer := range peers {
		if peer.IPAddress == ipAddress {
			t.Errorf("Peer with IP address %s should not exist", ipAddress)
			return
		}
	}
}