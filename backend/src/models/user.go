package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	FirstName string             `json:"first_name" bson:"first_name"`
	Surname   string             `json:"surname" bson:"surname"`
	Email     string             `json:"email" bson:"email"`
	Birthdate string             `json:"birthdate" bson:"birthdate"`
	Password  string             `json:"password" bson:"password"`
	FileNames []string           `json:"file_names" bson:"file_names"`
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
