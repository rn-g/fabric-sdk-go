package cop

import (
	"testing"
)

func TestEnroll(t *testing.T) {
	fcs, err := NewFabricCOPServices("http://localhost:8888", "../test_resources/")
	if err != nil {
		t.Errorf("NewFabricCOPServices return error: %v", err)
		return
	}
	_, err = fcs.Enroll("admin", "adminpw")
	if err != nil {
		t.Errorf("Enroll return error: %v", err)
		return
	}
}
