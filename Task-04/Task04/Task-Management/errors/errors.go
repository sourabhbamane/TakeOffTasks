package errors

//ServiceError is used to return business error messages
type CustomError struct {
	Message string `json:"message"`
}
