package logout

import (
	"encoding/json"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("LogoutUser", Logout)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "token",
		Value:  "", // Set an empty value
		MaxAge: -1, // Expire immediately
		Path:   "/",
	}

	http.SetCookie(w, cookie)
	json.NewEncoder(w).Encode(CustomError{Message: "Logged out successfully"})
}

type CustomError struct {
	Message string `json:"message"`
}
