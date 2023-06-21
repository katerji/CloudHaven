package service

import (
	"errors"
	"github.com/katerji/UserAuthKit/db"
	"github.com/katerji/UserAuthKit/db/query"
	"github.com/katerji/UserAuthKit/db/queryrow"
	"github.com/katerji/UserAuthKit/input"
	"github.com/katerji/UserAuthKit/model"
	"golang.org/x/crypto/bcrypt"
)

type authService struct{}

func (service authService) Register(input input.AuthInput) (int, error) {
	password, err := hashPassword(input.Password)
	if err != nil {
		return 0, err
	}
	input.Password = password
	return db.GetDbInstance().Insert(query.InsertUserQuery, input.Email, input.Name, input.Password)
}

func (service authService) Login(input input.AuthInput) (model.User, error) {
	result := queryrow.UserQueryRow{}
	client := db.GetDbInstance()
	row := client.QueryRow(query.GetUserByEmailQuery, input.Email)
	err := row.Scan(&result.ID, &result.Name, &result.Email, &result.Password)
	if err != nil {
		return model.User{}, errors.New("email does not exist")
	}

	if !validPassword(result.Password, input.Password) {
		return model.User{}, errors.New("incorrect password")
	}

	return model.User{
		ID:    result.ID,
		Name:  result.Name,
		Email: result.Email,
	}, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func validPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
