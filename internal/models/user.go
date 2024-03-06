package models

type UserInDatabase struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	FullName     string `json:"full_name,omitempty"`
}

type CreateUserRequest struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name,omitempty"`
}
