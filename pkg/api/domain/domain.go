package domain

type User struct {
	ID        int64  `json:"id,omitempty"`
	Firstname string `json:"firtsname" validated:"required"`
	Lastname  string `json:"lastname" validated:"required"`
	Password  string `json:"password,omitempty"`
	Email     string `json:"email" validated:"required"`
	CreatedOn int64  `json:"createdOn,omitempty"`
	Status    string `json:"status" validated:"required"`
}
