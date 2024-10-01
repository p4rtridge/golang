package api

import (
	"order_service/internal/core"
	"order_service/pkg"
	"order_service/services/product/entity"
	productUc "order_service/services/product/usecase"

	"github.com/gofiber/fiber/v2"
)

type ProductService interface {
	CreateProduct(*fiber.Ctx) error
	GetProducts(*fiber.Ctx) error
	SearchProducts(*fiber.Ctx) error
	GetProduct(*fiber.Ctx) error
	UpdateProduct(*fiber.Ctx) error
	DeleteProduct(*fiber.Ctx) error
}

type service struct {
	usecase productUc.ProductUsecase
}

func NewService(uc productUc.ProductUsecase) ProductService {
	return &service{
		usecase: uc,
	}
}

func (srv *service) CreateProduct(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	var data entity.ProductRequest
	if err := pkg.MultipartParser(form, &data); err != nil {
		return pkg.WriteResponse(c, core.ErrBadRequest)
	}
	if err := data.Validate(); err != nil {
		return pkg.WriteResponse(c, err)
	}

	err = srv.usecase.CreateProduct(c.Context(), &data)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(core.ResponseData(true))
}

func (srv *service) GetProducts(c *fiber.Ctx) error {
	products, err := srv.usecase.GetProducts(c.Context())
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(products))
}

func (srv *service) SearchProducts(c *fiber.Ctx) error {
	nameQuery := c.Query("name")
	if nameQuery == "" {
		return pkg.WriteResponse(c, core.ErrBadRequest)
	}

	products, err := srv.usecase.SearchProducts(c.Context(), nameQuery)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(products))
}

func (srv *service) GetProduct(c *fiber.Ctx) error {
	targetId, err := c.ParamsInt("productID")
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	product, err := srv.usecase.GetProduct(c.Context(), targetId)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(product))
}

func (srv *service) UpdateProduct(c *fiber.Ctx) error {
	targetId, err := c.ParamsInt("productID")
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	var data entity.ProductRequest
	if err := pkg.MultipartParser(form, &data); err != nil {
		return pkg.WriteResponse(c, core.ErrBadRequest)
	}
	if err := data.Validate(); err != nil {
		return pkg.WriteResponse(c, err)
	}

	err = srv.usecase.UpdateProduct(c.Context(), targetId, &data)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(true))
}

func (srv *service) DeleteProduct(c *fiber.Ctx) error {
	targetId, err := c.ParamsInt("productID")
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	err = srv.usecase.DeleteProduct(c.Context(), targetId)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusNoContent).JSON(core.ResponseData(true))
}
