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

// Create Product godoc
// @summary Create Product
// @description Create a new product with input payload
// @tags products
// @accept multipart/form-data
// @security BearerAuth
// @success 201
// @failure 400 {object} core.DefaultError
// @failure 401 {object} core.DefaultError
// @failure 404 {object} core.DefaultError
// @failure 500 {object} core.DefaultError
// @router /products/ [post]
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

	requester, ok := c.Locals(core.KeyRequester).(core.Requester)
	if !ok {
		return pkg.WriteResponse(c, core.ErrUnauthorized)
	}
	ctx := core.ContextWithRequester(c.Context(), requester)

	err = srv.usecase.CreateProduct(ctx, &data)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(core.ResponseData(true))
}

// Get Products godoc
// @summary Get Products
// @description Get entire products
// @tags products
// @success 200 {array} entity.Product
// @failure 404 {object} core.DefaultError
// @failure 500 {object} core.DefaultError
// @router /products/ [get]
func (srv *service) GetProducts(c *fiber.Ctx) error {
	products, err := srv.usecase.GetProducts(c.Context())
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(products))
}

// Search Products godoc
// @summary Search Products
// @description Search products with input query string
// @tags products
// @param name query string true "Query string"
// @success 200 {array} entity.Product
// @failure 400 {object} core.DefaultError
// @failure 404 {object} core.DefaultError
// @failure 500 {object} core.DefaultError
// @router /products/search [get]
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

// Get Product godoc
// @summary Get Product
// @description Get specific product
// @tags products
// @param productID path string true "Product's ID"
// @success 200 {object} entity.Product
// @failure 404 {object} core.DefaultError
// @failure 500 {object} core.DefaultError
// @router /products/:productID [get]
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

// Update Product godoc
// @summary Update Product
// @description Update the specific product
// @tags products
// @security BearerAuth
// @param productID path string true "Product's ID"
// @success 200
// @failure 400 {object} core.DefaultError
// @failure 401 {object} core.DefaultError
// @failure 404 {object} core.DefaultError
// @failure 500 {object} core.DefaultError
// @router /products/:productID [put]
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

	requester, ok := c.Locals(core.KeyRequester).(core.Requester)
	if !ok {
		return pkg.WriteResponse(c, core.ErrUnauthorized)
	}
	ctx := core.ContextWithRequester(c.Context(), requester)

	err = srv.usecase.UpdateProduct(ctx, targetId, &data)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(core.ResponseData(true))
}

// Delete Product godoc
// @summary Delete Product
// @description Delete the specific product
// @tags products
// @security BearerAuth
// @param productID path string true "Product's ID"
// @success 204
// @failure 400 {object} core.DefaultError
// @failure 401 {object} core.DefaultError
// @failure 500 {object} core.DefaultError
// @router /products/:productID [delete]
func (srv *service) DeleteProduct(c *fiber.Ctx) error {
	targetId, err := c.ParamsInt("productID")
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	requester, ok := c.Locals(core.KeyRequester).(core.Requester)
	if !ok {
		return pkg.WriteResponse(c, core.ErrUnauthorized)
	}
	ctx := core.ContextWithRequester(c.Context(), requester)

	err = srv.usecase.DeleteProduct(ctx, targetId)
	if err != nil {
		return pkg.WriteResponse(c, err)
	}

	return c.Status(fiber.StatusNoContent).JSON(core.ResponseData(true))
}
