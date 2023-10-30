package entity

type Employee struct {
	ID        int64  `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	PhoneNo   string `json:"phone"`
	Role      string `json:"role"`
	Salary    int64  `json:"salary"`
}
