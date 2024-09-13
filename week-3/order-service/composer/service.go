package composer

import (
	authAPI "order_service/services/auth/controller/api"
	userAPI "order_service/services/user/controller/api"

	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	Register(*fiber.Ctx) error
	Login(*fiber.Ctx) error
}

type UserService interface {
	GetUserProfile(*fiber.Ctx) error
}

func ComposeAuthAPIService(biz authAPI.AuthUseCase) AuthService {
	serviceAPI := authAPI.NewAPI(biz)

	return serviceAPI
}

func ComposeUserAPIService(biz userAPI.UserUsecase) UserService {
	serviceAPI := userAPI.NewAPI(biz)

	return serviceAPI
}
