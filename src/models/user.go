package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string   `json:"id"`
	FirstName string   `json:"first_name"`
	Surname   string   `json:"surname"`
	Email     string   `json:"email"`
	Birthdate string   `json:"birthdate"`
	Password  string   `json:"password"`
	FileNames []string `json:"file_names"`
}

// Hashes the user's password before storing it
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// Compares the user's password with the given password
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
