package entity

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/partridge1307/gofiber/pkg/validate"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" validate:"required,ascii,lowercase"`
	Password string `json:"-" validate:"required,gte=8"`
	Balance  float32
}

type UserValidationOutcome struct {
	Errors []string
	Valid  bool
}

type UserRepo interface {
	CreateUser(ctx context.Context, user *User) error
	GetUsers(ctx context.Context) (*[]User, error)
	GetUser(ctx context.Context, u interface{}) (*User, error)
}

func NewUser(id int, username, password string, balance float32) *User {
	return &User{
		ID:       id,
		Username: username,
		Password: password,
		Balance:  balance,
	}
}

func (u *User) SetID(id int) {
	u.ID = id
}

func (u *User) SetUsername(username string) {
	u.Username = username
}

func (u *User) SetPassword(password string) {
	u.Password = password
}

func (u *User) Validate(v *validator.Validate) *UserValidationOutcome {
	user := NewUser(u.ID, u.Username, u.Password, u.Balance)

	outcome := &UserValidationOutcome{
		Valid:  false,
		Errors: make([]string, 0),
	}

	if errs := v.Struct(user); errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			field := validate.ParseJsonField(User{}, err.StructField())

			outcome.Errors = append(outcome.Errors, fmt.Sprintf("%s needs to implement '%s'", field, err.Tag()))
		}
	}

	if len(outcome.Errors) == 0 {
		outcome.Valid = true
	}

	return outcome
}
