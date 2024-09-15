package composer

import (
	authAPI "order_service/services/auth/controller/api"
	productAPI "order_service/services/product/controller/api"
	userAPI "order_service/services/user/controller/api"

	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	Register(*fiber.Ctx) error
	Login(*fiber.Ctx) error
	Refresh(*fiber.Ctx) error
	SignOut(*fiber.Ctx) error
	SignOutAll(*fiber.Ctx) error
}

type UserService interface {
	GetUsers(*fiber.Ctx) error
	GetUser(*fiber.Ctx) error
	GetUserProfile(*fiber.Ctx) error
	AddUserBalance(*fiber.Ctx) error
}

type ProductService interface {
	CreateProduct(*fiber.Ctx) error
	GetProducts(*fiber.Ctx) error
	GetProduct(*fiber.Ctx) error
	UpdateProduct(*fiber.Ctx) error
	DeleteProduct(*fiber.Ctx) error
}

func ComposeAuthAPIService(biz authAPI.AuthUseCase) AuthService {
	serviceAPI := authAPI.NewAPI(biz)

	return serviceAPI
}

func ComposeUserAPIService(biz userAPI.UserUsecase) UserService {
	serviceAPI := userAPI.NewAPI(biz)

	return serviceAPI
}

func ComposeProductAPIService(biz productAPI.ProductUsecase) ProductService {
	serviceAPI := productAPI.NewAPI(biz)

	return serviceAPI
}
