package entity

type User struct {
	UserId   int64  `json:"userid"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
