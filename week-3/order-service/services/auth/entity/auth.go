package entity

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

func NewAuth(username, password string, role int) Auth {
	if role != -1 {
		return Auth{
			Username: username,
			Password: password,
			Role:     role,
		}
	}

	return Auth{
		Username: username,
		Password: password,
		Role:     0,
	}
}
