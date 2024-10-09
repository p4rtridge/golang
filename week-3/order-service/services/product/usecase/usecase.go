package usecase

import (
	"context"
	"order_service/internal/core"
	"order_service/services/product/entity"
	productAWSRepo "order_service/services/product/repository/aws"
	productPGRepo "order_service/services/product/repository/postgres"
)

type ProductUsecase interface {
	CreateProduct(ctx context.Context, data *entity.ProductRequest) error
	GetProducts(ctx context.Context) (*[]entity.Product, error)
	SearchProducts(ctx context.Context, nameQuery string) (*[]entity.Product, error)
	GetProduct(ctx context.Context, productID int) (*entity.Product, error)
	UpdateProduct(ctx context.Context, productID int, data *entity.ProductRequest) error
	DeleteProduct(ctx context.Context, productID int) error
}

type productUsecase struct {
	repo      productPGRepo.ProductRepository
	awsClient productAWSRepo.AWSClient
}

func NewUsecase(repo productPGRepo.ProductRepository, awsClient productAWSRepo.AWSClient) ProductUsecase {
	return &productUsecase{
		repo,
		awsClient,
	}
}

func (uc *productUsecase) CreateProduct(ctx context.Context, data *entity.ProductRequest) error {
	requester := core.GetRequester(ctx)
	uid, err := core.DecomposeUID(requester.GetSubject())
	if err != nil {
		return core.ErrInternalServerError.WithDebug(err.Error())
	}

	role := uid.GetRole()
	if role != 1 {
		return core.ErrBadRequest.WithError(entity.ErrCannotCreate.Error())
	}

	imageUrl, err := uc.awsClient.SaveImage(ctx, &data.Image)
	if err != nil {
		return core.ErrInternalServerError.WithError(entity.ErrCannotCreate.Error()).WithDebug(err.Error())
	}

	newProduct := entity.NewProduct(0, data.Name, imageUrl, data.Quantity, data.Price)

	err = uc.repo.CreateProduct(ctx, newProduct)
	if err != nil {
		return core.ErrInternalServerError.WithError(entity.ErrCannotCreate.Error()).WithDebug(err.Error())
	}

	return nil
}

func (uc *productUsecase) GetProducts(ctx context.Context) (*[]entity.Product, error) {
	products, err := uc.repo.GetProducts(ctx)
	if err != nil {
		if err == core.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}

		return nil, core.ErrInternalServerError.WithDebug(err.Error())
	}

	return products, nil
}

func (uc *productUsecase) SearchProducts(ctx context.Context, nameQuery string) (*[]entity.Product, error) {
	products, err := uc.repo.SearchProducts(ctx, nameQuery)
	if err != nil {
		if err == core.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}

		return nil, core.ErrInternalServerError.WithDebug(err.Error())
	}

	return products, nil
}

func (uc *productUsecase) GetProduct(ctx context.Context, productID int) (*entity.Product, error) {
	product, err := uc.repo.GetProduct(ctx, productID)
	if err != nil {
		if err == core.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}

		return nil, core.ErrInternalServerError.WithDebug(err.Error())
	}

	return product, nil
}

func (uc *productUsecase) UpdateProduct(ctx context.Context, productID int, data *entity.ProductRequest) error {
	requester := core.GetRequester(ctx)
	uid, err := core.DecomposeUID(requester.GetSubject())
	if err != nil {
		return core.ErrInternalServerError.WithDebug(err.Error())
	}

	role := uid.GetRole()
	if role != 1 {
		return core.ErrBadRequest.WithError(entity.ErrCannotUpdate.Error())
	}

	product, err := uc.repo.GetProduct(ctx, productID)
	if err != nil {
		return core.ErrNotFound.WithError(entity.ErrCannotUpdate.Error()).WithDebug(err.Error())
	}

	err = uc.awsClient.DeleteImage(ctx, product.ImageURL)
	if err != nil {
		return core.ErrInternalServerError.WithError(entity.ErrCannotUpdate.Error()).WithDebug(err.Error())
	}

	imageUrl, err := uc.awsClient.SaveImage(ctx, &data.Image)
	if err != nil {
		return core.ErrInternalServerError.WithError(entity.ErrCannotUpdate.Error()).WithDebug(err.Error())
	}

	updatedProduct := entity.NewProduct(productID, data.Name, imageUrl, data.Quantity, data.Price)
	err = uc.repo.UpdateProduct(ctx, productID, updatedProduct)
	if err != nil {
		return core.ErrInternalServerError.WithError(entity.ErrCannotUpdate.Error()).WithDebug(err.Error())
	}

	return nil
}

func (uc *productUsecase) DeleteProduct(ctx context.Context, productID int) error {
	requester := core.GetRequester(ctx)
	uid, err := core.DecomposeUID(requester.GetSubject())
	if err != nil {
		return core.ErrInternalServerError.WithDebug(err.Error())
	}

	role := uid.GetRole()
	if role != 1 {
		return core.ErrBadRequest.WithError(entity.ErrCannotDelete.Error())
	}

	err = uc.repo.DeleteProduct(ctx, productID)
	if err != nil {
		return core.ErrBadRequest.WithError(entity.ErrCannotDelete.Error()).WithDebug(err.Error())
	}

	return nil
}
