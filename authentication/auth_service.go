package authentication

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"

	"io"
	"io/ioutil"
)

/*Login is the HTTP Server response for calls to Login a user.*/
func Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		fmt.Println("Couldn't read request")
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		fmt.Println("Couldn't close channel for login request")
		panic(err)
	}

	var authInfo ClientUser
	if err := json.Unmarshal(body, &authInfo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		panic(err)
	}
	sessionToken, err := loginUser(authInfo)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Session could not be created, login info invalid"))
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(30 * time.Minute),
	})
}

/*Register is the HTTP server response for calls to Register a user.*/
func Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		fmt.Println("Couldn't read request")
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		fmt.Println("Couldn't close channel for login request")
		panic(err)
	}

	var authInfo ClientUser
	if err := json.Unmarshal(body, &authInfo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		panic(err)
	}

	err = registerUser(authInfo)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Registration failed"))
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Registration successful"))
}

func registerUser(authInfo ClientUser) error {
	err := insertUser(authInfo)
	return err
}

func loginUser(authInfo ClientUser) (string, error) {
	serverAuthInfo, err := getServerUserByUsername(authInfo.Username)
	if comparePasswords(serverAuthInfo.PasswordHash, []byte(authInfo.Password)) {
		//WE SHOULD FIRST CHECK IF THE USER HAS A COOKIE AND THEN CHECK IF SESSIONKEY HAS VALID LASTSEENTIME
		//IF TIME INVALID, COMPARE PASSWORD AND ISSUE A NEW SESSIONKEY
		//IF TIME VALID UPDATE SESSIONKEY IN USERSESSION TABLE, LET EM GO (MAYBE ISSUE NEW COOKIE?)

		userID := serverAuthInfo.Userid
		sessionToken := uuid.New().String()
		currentTime := time.Now().Format(time.RFC3339)
		expireTime := time.Now().Add(time.Minute * 30).Format(time.RFC3339)
		err = insertSessionToken(userID, sessionToken, currentTime, expireTime)
		if err != nil {
			fmt.Println("Unable to create session cookie")
			return "", err
		}
		//User is verified, return a cookie
		return sessionToken, err
	}
	return "", err
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
