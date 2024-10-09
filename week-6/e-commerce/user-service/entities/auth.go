package entities

type Role int

const (
	USER Role = iota
	ADMIN
)

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

func NewAuth(username, password string, role Role) Auth {
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
		Role:     USER,
	}
}
