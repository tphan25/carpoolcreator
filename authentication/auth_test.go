package authentication

import (
	"fmt"
	"testing"
)

func TestHashing(t *testing.T) {
	//Let's do a quick test then
	testPass := "yesThisIsPass"
	testPass = hashAndSalt([]byte(testPass))
	fmt.Println("This is salted pass", testPass)

	testPass2 := "noThisNotPass"
	testPass2 = hashAndSalt([]byte(testPass2))
	fmt.Println("This is salted not pass", testPass2)

	if !comparePasswords(testPass, []byte("noThisIsntPass")) {
		t.Errorf(testPass)
	}
}
