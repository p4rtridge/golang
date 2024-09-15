package composer

import (
	"order_service/config"
	"order_service/pkg"
	authAPI "order_service/services/auth/controller/api"
	authPGRepo "order_service/services/auth/repository/postgres"
	authRDRepo "order_service/services/auth/repository/redis"
	authUsecase "order_service/services/auth/usecase"
	userAPI "order_service/services/user/controller/api"
	userPGRepo "order_service/services/user/repository/postgres"
	userUsecase "order_service/services/user/usecase"
	"runtime"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func ComposeAuthBusiness(cfg *config.Config, pg *pgxpool.Pool, rd *redis.Client) authAPI.AuthUseCase {
	repo := authPGRepo.NewPostgresRepo(pg)
	tokenRepo := authRDRepo.NewRedisRepo(rd)
	hasher := pkg.NewHasher(64*1024, 3, 16, 32, uint8(runtime.NumCPU()))
	jwt := pkg.NewJWT(cfg.SecretKey, cfg.ATExpireInSec, cfg.RTExpireInSec)

	return authUsecase.NewUsecase(repo, tokenRepo, hasher, jwt)
}

func ComposeUserBusiness(cfg *config.Config, db *pgxpool.Pool) userAPI.UserUsecase {
	repo := userPGRepo.NewPostgresRepo(db)

	return userUsecase.NewUsecase(repo)
}
