package entity

import "order_service/pkg"

type AuthUsernamePassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *AuthUsernamePassword) Validate() error {
	if err := pkg.UsernameIsValid(a.Username); err != nil {
		return err
	}

	if err := pkg.CheckPassword(a.Password); err != nil {
		return err
	}

	return nil
}

type Token struct {
	Token     string `json:"token"`
	ExpiredIn int    `json:"expire_in"`
}

type TokenResponse struct {
	AccessToken Token `json:"access_token"`
}
