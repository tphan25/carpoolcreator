package authentication

import (
	"testing"
)

/*TestHashing tests if our hash/salt function for passwords functions. */
func TestHashing(t *testing.T) {
	testPass := "yesThisIsPass"
	testPass = hashAndSalt([]byte(testPass))

	testPass2 := "noThisNotPass"
	testPass2 = hashAndSalt([]byte(testPass2))

	if !comparePasswords(testPass, []byte("yesThisIsPass")) {
		t.Errorf(testPass)
		t.Errorf("Password did not match")
	}
}

/*TestRegister tests if a user can register a username and password.
1. Call auth_service.go:registerUser with test credentials
2. Call auth_service.go:loginUser with same credentials
3. TODO: Delete the user afterwards
*/
func TestRegister(t *testing.T) {
	err := Init()
	if err != nil {
		t.Errorf("Error occurred when creating session")
		t.Errorf(err.Error())
		t.FailNow()
	}

	authInfo := ClientUser{
		Username:  "Test User",
		Password:  "Test Pass",
		FirstName: "TestFirst",
		LastName:  "TestLast",
	}

	err = deleteUser(authInfo.Username)
	if err != nil {
		t.Errorf("Error occurred when deleting user")
		t.Errorf(err.Error())
		t.FailNow()
	}

	err = registerUser(authInfo)
	if err != nil {
		t.Errorf("Error occurred when registering user")
		t.Errorf(err.Error())
		t.FailNow()
	}

	_, err = loginUser(authInfo)
	if err != nil {
		t.Errorf("Error occurred when logging in user")
		t.Errorf(err.Error())
		t.FailNow()
	}
}
