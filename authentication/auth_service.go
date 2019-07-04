package authentication

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"io"
	"io/ioutil"
)

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		fmt.Println("Couldn't read request")
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		fmt.Println("Couldn't close channel for login request")
	}

	var authInfo ClientUser
	if err := json.Unmarshal(body, &authInfo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	serverAuthInfo, err := getServerUserByUsername(authInfo.Username)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	if comparePasswords(serverAuthInfo.PasswordHash, []byte(authInfo.Password)) {
		//User is verified, return a token
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hey password was same"))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Hey it wasn't the same"))
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		fmt.Println("Couldn't read request")
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		fmt.Println("Couldn't close channel for login request")
	}

	var authInfo RegisterUser
	if err := json.Unmarshal(body, &authInfo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		fmt.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	err = insertUser(authInfo)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("yeah u aint allowed"))
		fmt.Println(err)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("yeah u allowed"))
	}
}

func comparePasswords(hashedPass string, plainPass []byte) bool {
	byteHash := []byte(hashedPass)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPass)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func hashAndSalt(pwd []byte) string {
	//Second field represents cost of hashing, higher means more CPU power to hash
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
