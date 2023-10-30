package errors

//ServiceError is used to return business error messages
type ServiceError struct {
	Message string `json:"message"`
}
