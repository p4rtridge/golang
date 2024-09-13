package composer

import (
	"order_service/config"
	"order_service/pkg"
	authAPI "order_service/services/auth/controller/api"
	authPGRepo "order_service/services/auth/repository/postgres"
	authUsecase "order_service/services/auth/usecase"
	userAPI "order_service/services/user/controller/api"
	userPGRepo "order_service/services/user/repository/postgres"
	userUsecase "order_service/services/user/usecase"
	"runtime"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ComposeAuthBusiness(cfg *config.Config, db *pgxpool.Pool) authAPI.AuthUseCase {
	repo := authPGRepo.NewPostgresRepo(db)
	hasher := pkg.NewHasher(64*1024, 3, 16, 32, uint8(runtime.NumCPU()))
	jwt := pkg.NewJWT(cfg.SecretKey, cfg.ExpireTokenInSec)

	return authUsecase.NewUsecase(repo, hasher, jwt)
}

func ComposeUserBusiness(cfg *config.Config, db *pgxpool.Pool) userAPI.UserUsecase {
	repo := userPGRepo.NewPostgresRepo(db)

	return userUsecase.NewUsecase(repo)
}
