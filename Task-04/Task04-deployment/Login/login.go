package login

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/api/iterator"
)

func init() {
	functions.HTTP("LoginUser", LoginHandler)
}

type User struct {
	UserId   int64  `json:"userid"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// declaring Constants
const (
	projectId string = "task-management-405310"
)

var jwtKey = []byte("secret_key")

// this struct will be using for payload for jwt
// inside payload we will pass username,role and when the token is expiring
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
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,PATCH,DELETE,OPTIONS")
	w.Header().Set("Content-Type", "application/json")

	var credentials Credentials

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// here we will Validate username and password
	_, err := ValidateUser(credentials.Username, credentials.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(CustomError{Message: "Invalid credentials"})
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
	json.NewEncoder(w).Encode(CustomError{Message: "Logged in Succesfully"})
}

func ValidateUser(username, password string) (*User, error) {
	ctx := context.Background()
	//creating new client on firestore
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		//log.Write("Failed to create Firestore Client")
		return nil, err
	}
	defer client.Close()
	// Reference to the Firestore collection where employees are stored.
	collection := client.Collection("users")
	// Query Firestore to get the user with the provided username and password
	iter := collection.Where("Username", "==", username).Where("Password", "==", password).Documents(ctx)
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, errors.New("user not found")
	}
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return nil, err
	}

	//var uname = username
	var user User
	err = doc.DataTo(&user)
	if err != nil {
		log.Printf("Error converting user data: %v", err)
		return nil, err
	}

	return &user, nil
}

// ServiceError is used to return business error messages
type CustomError struct {
	Message string `json:"message"`
}
