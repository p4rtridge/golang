package composer

import (
	"order_service/config"
	"order_service/pkg"
	authPGRepo "order_service/services/auth/repository/postgres"
	authRDRepo "order_service/services/auth/repository/redis"
	authUsecase "order_service/services/auth/usecase"
	orderPGRepo "order_service/services/order/repository/postgres"
	orderUsecase "order_service/services/order/usecase"
	productPGRepo "order_service/services/product/repository/postgres"
	productUsecase "order_service/services/product/usecase"
	userPGRepo "order_service/services/user/repository/postgres"
	userUsecase "order_service/services/user/usecase"
	"runtime"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func ComposeAuthUsecase(cfg *config.Config, pg *pgxpool.Pool, rd *redis.Client) authUsecase.AuthUseCase {
	repo := authPGRepo.NewAuthRepo(pg)
	tokenRepo := authRDRepo.NewTokenRepo(rd)
	hasher := pkg.NewHasher(64*1024, 3, 16, 32, uint8(runtime.NumCPU()))
	jwt := pkg.NewJWT(cfg.SecretKey, cfg.ATExpireInSec, cfg.RTExpireInSec)

	return authUsecase.NewUsecase(repo, tokenRepo, hasher, jwt)
}

func ComposeUserUsecase(db *pgxpool.Pool) userUsecase.UserUsecase {
	repo := userPGRepo.NewUserRepo(db)

	return userUsecase.NewUsecase(repo)
}

func ComposeProductUsecase(db *pgxpool.Pool) productUsecase.ProductUsecase {
	repo := productPGRepo.NewProductRepo(db)

	return productUsecase.NewUsecase(repo)
}

func ComposeOrderUsecase(db *pgxpool.Pool) orderUsecase.OrderUsecase {
	repo := orderPGRepo.NewOrderRepo(db)

	return orderUsecase.NewUsecase(repo)
}
