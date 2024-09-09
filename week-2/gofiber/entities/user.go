package entities

type User struct {
	ID       int     `json:"id"`
	Username string  `json:"user_name"`
	Password string  `json:"password,omitempty"`
	Balance  float32 `json:"balance"`
}

func (u *User) SetPassword(password string) {
	u.Password = password
}
