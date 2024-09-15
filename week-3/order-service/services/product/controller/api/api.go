package api

import (
	"context"
	"order_service/internal/core"
	"order_service/pkg"
	"order_service/services/product/entity"

	"github.com/gofiber/fiber/v2"
)

type ProductUsecase interface {
	CreateProduct(ctx context.Context, data *entity.ProductRequest) error
	GetProducts(ctx context.Context) (*[]entity.Product, error)
	GetProduct(ctx context.Context, productID int) (*entity.Product, error)
	UpdateProduct(ctx context.Context, productID int, data *entity.ProductRequest) error
	DeleteProduct(ctx context.Context, productID int) error
}

type api struct {
	usecase ProductUsecase
}

func NewAPI(uc ProductUsecase) *api {
	return &api{
		usecase: uc,
	}
}

func (api *api) CreateProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var data entity.ProductRequest

	if err := c.BodyParser(&data); err != nil {
		return pkg.WriteResponse(c, err)
	}

	err := api.usecase.CreateProduct(ctx, &data)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(core.ResponseData(true))
}

func (api *api) GetProducts(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	products, err := api.usecase.GetProducts(ctx)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(products))
}

func (api *api) GetProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	targetId, err := c.ParamsInt("productID")
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	product, err := api.usecase.GetProduct(ctx, targetId)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(product))
}

func (api *api) UpdateProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	targetId, err := c.ParamsInt("productID")
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	var data entity.ProductRequest

	if err := c.BodyParser(&data); err != nil {
		return pkg.WriteResponse(c, err)
	}

	err = api.usecase.UpdateProduct(ctx, targetId, &data)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(true))
}

func (api *api) DeleteProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	targetId, err := c.ParamsInt("productID")
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	err = api.usecase.DeleteProduct(ctx, targetId)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusNoContent).JSON(core.ResponseData(true))
}
