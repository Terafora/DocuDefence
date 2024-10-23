package models

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
	Birthdate string `json:"birthdate"`
	FileName  string `json:"file_name"`
}
