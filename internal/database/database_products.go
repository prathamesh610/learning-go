package database

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/prathameshj610/go-microservices/internal/dberrors"
	"github.com/prathameshj610/go-microservices/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c Client) GetAllProducts(ctx context.Context, vendorId string) ([]models.Product, error) {
	var vendors []models.Product
	result := c.DB.WithContext(ctx).Where(models.Product{VendorId: vendorId}).
		Find(&vendors)

	return vendors, result.Error
}

func (c Client) AddProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	product.ProductId = uuid.NewString()
	result := c.DB.WithContext(ctx).Create(&product)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		}
		return nil, result.Error
	}

	return product, nil
}

func (c Client) GetProductById(ctx context.Context, ID string) (*models.Product, error) {
	product := &models.Product{}

	result := c.DB.WithContext(ctx).
				Where(&models.Product{ProductId: ID}).
				First(&product)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound){
			return nil, &dberrors.NotFoundError{Entity: "product", ID: ID}
		}
		return nil, result.Error
	}

	return product, nil
}

func (c Client) UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	var products []models.Product

	result := c.DB.WithContext(ctx).
				Model(&products).
				Clauses(clause.Returning{}).
				Where(&models.Product{ProductId:  product.ProductId}).
				Updates(models.Product{
					Name: product.Name,
					VendorId: product.VendorId,
					Price: product.Price,
				})

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey){
			return nil, &dberrors.ConflictError{}
		}
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, &dberrors.NotFoundError{Entity: "product", ID: product.ProductId}
	}

	return &products[0], nil

}	


func (c Client) DeleteProduct(ctx context.Context, ID string) error {
	return c.DB.WithContext(ctx).Delete(&models.Product{ProductId: ID,}).Error
}