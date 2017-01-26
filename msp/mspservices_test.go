package msp

import (
	"testing"
)

func TestEnrollWithMissingParameters(t *testing.T) {
	msps, err := NewMSPServices("localhost", "/")
	if err != nil {
		t.Fatalf("NewMSPServices return error: %v", err)
	}
	_, _, err = msps.Enroll("", "user1")
	if err == nil {
		t.Fatalf("Enroll didn't return error")
	}
	if err.Error() != "enrollmentID is empty" {
		t.Fatalf("Enroll didn't return right error")
	}
	_, _, err = msps.Enroll("test", "")
	if err == nil {
		t.Fatalf("Enroll didn't return error")
	}
	if err.Error() != "enrollmentSecret is empty" {
		t.Fatalf("Enroll didn't return right error")
	}
}
