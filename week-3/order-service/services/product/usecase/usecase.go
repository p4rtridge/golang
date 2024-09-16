package usecase

import (
	"context"
	"order_service/internal/core"
	"order_service/services/product/entity"
	productRepo "order_service/services/product/repository/postgres"
)

type ProductUsecase interface {
	CreateProduct(ctx context.Context, data *entity.ProductRequest) error
	GetProducts(ctx context.Context) (*[]entity.Product, error)
	GetProduct(ctx context.Context, productID int) (*entity.Product, error)
	UpdateProduct(ctx context.Context, productID int, data *entity.ProductRequest) error
	DeleteProduct(ctx context.Context, productID int) error
}

type productUsecase struct {
	repo productRepo.ProductRepository
}

func NewUsecase(repo productRepo.ProductRepository) ProductUsecase {
	return &productUsecase{
		repo,
	}
}

func (uc *productUsecase) CreateProduct(ctx context.Context, data *entity.ProductRequest) error {
	newProduct := entity.NewProduct(0, data.Name, data.Quantity, data.Price)

	err := uc.repo.CreateProduct(ctx, &newProduct)
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
	newData := entity.NewProduct(0, data.Name, data.Quantity, data.Price)

	err := uc.repo.UpdateProduct(ctx, productID, &newData)
	if err != nil {
		return core.ErrInternalServerError.WithError(entity.ErrCannotUpdate.Error()).WithDebug(err.Error())
	}

	return nil
}

func (uc *productUsecase) DeleteProduct(ctx context.Context, productID int) error {
	err := uc.repo.DeleteProduct(ctx, productID)
	if err != nil {
		return core.ErrBadRequest.WithError(entity.ErrCannotDelete.Error()).WithDebug(err.Error())
	}

	return nil
}
