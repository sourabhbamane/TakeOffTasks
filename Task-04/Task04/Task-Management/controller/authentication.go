package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/task-management/errors"
	"github.com/task-management/repository"
)

var jwtKey = []byte("secret_key")

// this struct will be using for payload for jwt
// inside payload we will pass username and when the token is expiring
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// to create token
func createToken(username string) (string, error) {
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	//in login we will check if the user is valid or not
	//if valid create claim object and after creating obj create jwt token from it
	//and store that data into cookies and pass the data into it

	var credentials Credentials

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// here we will Validate username and password
	_, err := repository.NewFireStoreRepo().ValidateUser(credentials.Username, credentials.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errors.CustomError{Message: "Invalid credentials"})
		return
	}

	//to create token
	tokenString, err := createToken(credentials.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//TO Set the cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: time.Now().Add(30 * time.Minute),
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(errors.CustomError{Message: "Logged in Succesfully"})
}

// To LogOut the LoggedIn User
// Written this Cause we have set the token expire time as 30 min
// to logout imiadiatally written this functionality
func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "token",
		Value:  "", // Set an empty value
		MaxAge: -1, // Expire immediately
		Path:   "/",
	}

	http.SetCookie(w, cookie)
	json.NewEncoder(w).Encode(errors.CustomError{Message: "Logged out successfully"})
}
