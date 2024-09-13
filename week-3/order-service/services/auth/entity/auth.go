package entity

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewAuth(username, password string) Auth {
	return Auth{
		Username: username,
		Password: password,
	}
}
