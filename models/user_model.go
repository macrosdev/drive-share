package models

import (
	"context"
	"errors"
	"server/utilities"

	emailverifier "github.com/AfterShip/email-verifier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required,min=8"`
	Avatar    string `json:"avatar"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var verifier = emailverifier.NewVerifier().EnableSMTPCheck()

type User struct {
	Id        primitive.ObjectID `json:"id" binding:"required"`
	Firstname string             `json:"firstname" binding:"required"`
	Lastname  string             `json:"lastname" binding:"required"`
	Username  string             `json:"username" binding:"required"`
	Email     string             `json:"email" binding:"required,email"`
	Password  string             `json:"password" binding:"required,min=8"`
	Avatar    string             `json:"avatar" binding:"required"`
}

func (newUser *User) SaveUser(c context.Context) (*User, error) {

	_, err := userCollection.InsertOne(c, newUser)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("user with that email already exists")
		}
	}

	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: opt}

	if _, err := userCollection.Indexes().CreateOne(c, index); err != nil {
		return nil, errors.New("could not create index of email")
	}

	return nil, nil
}

func ValidateEmail(email string) (bool, error) {
	_, err := verifier.Verify(email)

	if err != nil {
		return false, errors.New("not registered email")
	}

	return true, nil
}

func IsEmailRegistered(email string, c context.Context) bool {
	user := User{}
	err := userCollection.FindOne(c, bson.M{"email": email}).Decode(&user)
	return err == nil
}

func LoginCheck(email string, password string, c context.Context) (User, error) {
	var err error

	user := User{}
	err = userCollection.FindOne(c, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return User{}, err
	}

	err = utilities.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return User{}, err
	}

	return user, nil
}
