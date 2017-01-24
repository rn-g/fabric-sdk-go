package cop

import (
	"testing"
)

func TestEnroll(t *testing.T) {
	msps, err := NewMSPServices("http://localhost:8888", "./test_resources/")
	if err != nil {
		t.Fatalf("NewMSPServices return error: %v", err)
	}
	key, cert, err := msps.Enroll("admin", "adminpw")
	if err != nil {
		t.Fatalf("Enroll return error: %v", err)
	}
	if key == nil {
		t.Fatalf("private key return from Enroll is nil")
	}
	if cert == nil {
		t.Fatalf("cert return from Enroll is nil")

	}
}
