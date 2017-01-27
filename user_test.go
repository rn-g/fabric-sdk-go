package fabric_sdk_go

import (
	"testing"
)

func TestUserMethods(t *testing.T) {
	user := NewUser("testUser")
	if user.GetName() != "testUser" {
		t.Fatalf("NewUser create wrong user")
	}
	var roles []string
	roles = append(roles, "admin")
	roles = append(roles, "user")
	user.SetRoles(roles)

	if user.GetRoles()[0] != "admin" {
		t.Fatalf("user.GetRoles() return wrong user")
	}
	if user.GetRoles()[1] != "user" {
		t.Fatalf("user.GetRoles() return wrong user")
	}

}
