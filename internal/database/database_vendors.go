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

func (c Client) GetAllVendors(ctx context.Context) ([]models.Vendor, error) {
	var vendors []models.Vendor
	result := c.DB.WithContext(ctx).Find(&vendors)

	return vendors, result.Error
}


func (c Client) AddVendor(ctx context.Context, vendor *models.Vendor) (*models.Vendor, error) {
	vendor.VendorId = uuid.NewString()
	result := c.DB.WithContext(ctx).Create(&vendor)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		}
		return nil, result.Error
	}

	return vendor, nil
}

func (c Client) GetVendorById(ctx context.Context, ID string) (*models.Vendor, error) {
	vendor := &models.Vendor{}

	result := c.DB.WithContext(ctx).
				Where(&models.Vendor{VendorId: ID}).
				First(&vendor)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &dberrors.NotFoundError{Entity: "vendor", ID: ID}
		}
		return nil, result.Error
	}
	return vendor, nil
}

func (c Client) UpdateVendor(ctx context.Context, vendor *models.Vendor) (*models.Vendor, error) {
	var vendors []models.Vendor

	result := c.DB.WithContext(ctx).
				Model(&vendors).
				Clauses(clause.Returning{}).
				Where(&models.Vendor{VendorId: vendor.VendorId}).
				Updates(models.Vendor{
					Name: vendor.Name,
					Contact: vendor.Contact,
					Email: vendor.Email,
					Address: vendor.Address,
				})

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		}
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, &dberrors.NotFoundError{Entity: "vendor", ID: vendor.VendorId}
	}

	return &vendors[0], nil

}

func (c Client) DeleteVendor(ctx context.Context, ID string) error {
	return c.DB.WithContext(ctx).Delete(&models.Vendor{VendorId: ID,}).Error
}