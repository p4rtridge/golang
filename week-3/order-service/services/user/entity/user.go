package entity

import "time"

type User struct {
	Id        int        `json:"id"`
	Username  string     `json:"username"`
	Password  string     `json:"-"` // sensitive field, should not send to user
	Role      int        `json:"role"`
	Balance   float32    `json:"balance"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func NewUser(id int, username, password string) User {
	return User{
		Id:        id,
		Username:  username,
		Password:  password,
		Balance:   0.0,
		CreatedAt: time.Now(),
		UpdatedAt: &time.Time{},
	}
}

func (user *User) SetBalance(b float32) {
	if user != nil {
		user.Balance = b
	}
}

func (user User) GetId() int {
	return user.Id
}

func (user User) GetBalance() float32 {
	return user.Balance
}
