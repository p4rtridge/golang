package test

import (
	"context"
	"errors"
	"order_service/internal/core"
	"order_service/services/product/entity"
	"order_service/services/product/test/mock"
	"order_service/services/product/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ProductUsecaseTestSuite struct {
	suite.Suite

	products    *[]entity.Product
	mockRepo    *mock.MockProductRepository
	mockAWSRepo *mock.MockAWSClient
	usecase     usecase.ProductUsecase
}

func (suite *ProductUsecaseTestSuite) SetupTest() {
	suite.products = &[]entity.Product{
		{
			Id:        1,
			Name:      "orange",
			Quantity:  10,
			Price:     2.5,
			CreatedAt: time.Now(),
		},
		{
			Id:        2,
			Name:      "pineapple",
			Quantity:  25,
			Price:     5,
			CreatedAt: time.Now(),
		},
	}

	ctrl := gomock.NewController(suite.T())

	suite.mockRepo = mock.NewMockProductRepository(ctrl)
	suite.mockAWSRepo = mock.NewMockAWSClient(ctrl)
	suite.usecase = usecase.NewUsecase(suite.mockRepo, suite.mockAWSRepo)
}

func (suite *ProductUsecaseTestSuite) TestCreateProduct() {
	tests := []struct {
		name      string
		ctx       context.Context
		data      *entity.ProductRequest
		repoErr   error
		want      error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Valid product",
			ctx:  context.TODO(),
			data: &entity.ProductRequest{
				Name:     "orange",
				Quantity: 25,
				Price:    2.5,
			},
			repoErr:   nil,
			want:      nil,
			assertion: assert.NoError,
		},
		{
			name:      "Invalid product",
			ctx:       context.TODO(),
			data:      &entity.ProductRequest{},
			repoErr:   errors.New("an error occur while creating product"),
			want:      core.ErrInternalServerError.WithError(entity.ErrCannotCreate.Error()).WithDebug(errors.New("an error occur while creating product").Error()),
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.mockRepo.EXPECT().CreateProduct(gomock.Any(), tt.data).Return(tt.repoErr)

			err := suite.usecase.CreateProduct(tt.ctx, tt.data)

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.want, "error should be return correctly")
			}
		})
	}
}

func (suite *ProductUsecaseTestSuite) TestGetProducts() {
	tests := []struct {
		name      string
		repoErr   error
		want      *[]entity.Product
		wantErr   error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Products exists",
			repoErr:   nil,
			want:      suite.products,
			wantErr:   nil,
			assertion: assert.NoError,
		},
		{
			name:      "Products empty",
			repoErr:   core.ErrRecordNotFound,
			want:      nil,
			wantErr:   core.ErrNotFound,
			assertion: assert.Error,
		},
		{
			name:      "Products return an error",
			repoErr:   errors.New("this is an error"),
			want:      nil,
			wantErr:   core.ErrInternalServerError.WithDebug(errors.New("this is an error").Error()),
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.mockRepo.EXPECT().GetProducts(gomock.Any()).Return(tt.want, tt.repoErr)

			products, err := suite.usecase.GetProducts(context.Background())

			suite.Equal(tt.want, products, "products should be retrieved correctly")

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.wantErr, "error should be return correctly")
			}
		})
	}
}

func (suite *ProductUsecaseTestSuite) TestGetProduct() {
	tests := []struct {
		name      string
		productId int
		repoErr   error
		want      *entity.Product
		wantErr   error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Product exists",
			productId: 1,
			repoErr:   nil,
			want:      &(*suite.products)[0],
			wantErr:   nil,
			assertion: assert.NoError,
		},
		{
			name:      "Product empty",
			productId: 1,
			repoErr:   core.ErrRecordNotFound,
			want:      nil,
			wantErr:   core.ErrNotFound,
			assertion: assert.Error,
		},
		{
			name:      "Repo return an error",
			productId: 1,
			repoErr:   errors.New("this is an error"),
			want:      nil,
			wantErr:   core.ErrInternalServerError.WithDebug(errors.New("this is an error").Error()),
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.mockRepo.EXPECT().GetProduct(gomock.Any(), tt.productId).Return(tt.want, tt.repoErr)

			products, err := suite.usecase.GetProduct(context.Background(), tt.productId)

			suite.Equal(tt.want, products, "product should be retrieved correctly")

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.wantErr, "error should be return correctly")
			}
		})
	}
}

func (suite *ProductUsecaseTestSuite) TestUpdateProduct() {
	tests := []struct {
		name      string
		productId int
		data      *entity.ProductRequest
		repoErr   error
		want      error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Valid product",
			productId: 1,
			data: &entity.ProductRequest{
				Name:     "apple",
				Quantity: 25,
				Price:    5.5,
			},
			repoErr:   nil,
			want:      nil,
			assertion: assert.NoError,
		},
		{
			name:      "Invalid product",
			productId: 1,
			data:      &entity.ProductRequest{},
			repoErr:   core.ErrRecordNotFound,
			want:      core.ErrInternalServerError.WithError(entity.ErrCannotUpdate.Error()).WithDebug(core.ErrRecordNotFound.Error()),
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.mockRepo.EXPECT().UpdateProduct(gomock.Any(), tt.productId, tt.data).Return(tt.repoErr)

			err := suite.usecase.UpdateProduct(context.Background(), tt.productId, tt.data)

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.want, "error should be return correctly")
			}
		})
	}
}

func (suite *ProductUsecaseTestSuite) TestDeleteProduct() {
	tests := []struct {
		name      string
		productId int
		repoErr   error
		want      error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Product exists",
			productId: 1,
			repoErr:   nil,
			want:      nil,
			assertion: assert.NoError,
		},
		{
			name:      "Product does not exists",
			productId: 1,
			repoErr:   core.ErrRecordNotFound,
			want:      core.ErrBadRequest.WithError(entity.ErrCannotDelete.Error()).WithDebug(core.ErrRecordNotFound.Error()),
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.mockRepo.EXPECT().DeleteProduct(gomock.Any(), tt.productId).Return(tt.repoErr)

			err := suite.usecase.DeleteProduct(context.Background(), tt.productId)

			if tt.assertion(suite.T(), err) {
				suite.ErrorIs(err, tt.want, "error should be return correctly")
			}
		})
	}
}

func TestProductUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(ProductUsecaseTestSuite))
}
