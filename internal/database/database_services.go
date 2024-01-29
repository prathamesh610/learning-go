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

func (c Client) GetAllServices(ctx context.Context) ([]models.Services, error){
	var services []models.Services
	result := c.DB.WithContext(ctx).Find(&services)

	return services, result.Error
}

func (c Client) AddService(ctx context.Context, service *models.Services) (*models.Services, error) {
	service.ServiceId = uuid.NewString()

	result := c.DB.WithContext(ctx).Create(&service)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey){
			return nil, &dberrors.ConflictError{}
		}
		return nil, result.Error
	}

	return service, nil
}

func (c Client) GetServiceById(ctx context.Context, ID string) (*models.Services, error) {
	service := &models.Services{}

	result := c.DB.WithContext(ctx).
				Where(&models.Services{ServiceId: ID}).
				First(&service)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound){
			return nil, &dberrors.NotFoundError{Entity: "service", ID: ID}
		}
		return nil, result.Error
	}

	return service, nil
}

func (c Client) UpdateService(ctx context.Context, service *models.Services) (*models.Services, error) {
	var services []models.Services

	result := c.DB.WithContext(ctx).
				Model(&services).
				Clauses(clause.Returning{}).
				Where(&models.Services{ServiceId: service.ServiceId}).
				Updates(models.Services{
					Name: service.Name,
					Price: service.Price,
				})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		}
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, &dberrors.NotFoundError{Entity: "service", ID: service.ServiceId}
	}

	return &services[0], nil
}

func (c Client) DeleteService(ctx context.Context, ID string) error {
	return c.DB.WithContext(ctx).Delete(&models.Services{ServiceId: ID,}).Error
}