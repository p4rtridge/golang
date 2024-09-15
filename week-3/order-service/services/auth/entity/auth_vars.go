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

type AuthLogin struct {
	AuthUsernamePassword
	DeviceId string `json:"device_id"`
}

func (a *AuthLogin) Validate() error {
	if err := a.AuthUsernamePassword.Validate(); err != nil {
		return err
	}

	if err := pkg.DeviceIdIsValid(a.DeviceId); err != nil {
		return err
	}

	return nil
}

type AuthSignOut struct {
	DeviceId string `json:"device_id"`
}

func (a *AuthSignOut) Validate() error {
	if err := pkg.DeviceIdIsValid(a.DeviceId); err != nil {
		return err
	}

	return nil
}

type Token struct {
	Token     string `json:"token,omitempty"`
	ExpiredIn int    `json:"expire_in,omitempty"`
}

type TokenResponse struct {
	AccessToken  Token `json:"access_token,omitempty"`
	RefreshToken Token `json:"refresh_token,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	DeviceId     string `json:"device_id"`
}
