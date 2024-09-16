package composer

import (
	authSrv "order_service/services/auth/controller/api"
	authUc "order_service/services/auth/usecase"
	orderSrv "order_service/services/order/controller/api"
	orderUc "order_service/services/order/usecase"
	productSrv "order_service/services/product/controller/api"
	productUc "order_service/services/product/usecase"
	userSrv "order_service/services/user/controller/api"
	userUc "order_service/services/user/usecase"
)

func ComposeAuthAPIService(biz authUc.AuthUseCase) authSrv.AuthService {
	serviceAPI := authSrv.NewService(biz)

	return serviceAPI
}

func ComposeUserAPIService(biz userUc.UserUsecase) userSrv.UserService {
	serviceAPI := userSrv.NewService(biz)

	return serviceAPI
}

func ComposeProductAPIService(biz productUc.ProductUsecase) productSrv.ProductService {
	serviceAPI := productSrv.NewService(biz)

	return serviceAPI
}

func ComposeOrderAPIService(biz orderUc.OrderUsecase) orderSrv.OrderService {
	serviceAPI := orderSrv.NewService(biz)

	return serviceAPI
}
