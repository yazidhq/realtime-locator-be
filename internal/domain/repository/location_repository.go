package repository

import (
	"TeamTrackerBE/internal/domain/model"
	"TeamTrackerBE/internal/utils"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LocationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) Create(location *model.Location) (*model.Location, error) {
	err := r.db.
		Create(location).
		Error

	return location, err
}

func (r *LocationRepository) BulkCreate(locations []model.Location) error {
	if len(locations) == 0 {
		return nil
	}
	
	return r.db.CreateInBatches(&locations, 200).Error
}

func (r *LocationRepository) Update(locationID uuid.UUID, req model.Location) (*model.Location, error) {
    var location model.Location

    if err := r.db.
		First(&location, locationID).
		Error; err != nil {
			return nil, err
		}

    if err := r.db.
		Model(&location).
		Updates(req).
		Error; err != nil {
			return nil, err
		}

    return &location, nil
}

func (r *LocationRepository) Delete(locationID uuid.UUID) (*model.Location, error) {
	var location model.Location

	err := r.db.
		Where("id = ?", locationID).
		Unscoped().
		Delete(&location).
		Error

	return &location, err
}

func (r *LocationRepository) FindAll(page, limit int, filters []utils.FilterOptions, sorts []utils.SortOption) ([]model.Location, int, error) {
    var locations []model.Location
    var total int64

    db := r.db.
		Model(&model.Location{})

    db = utils.ApplyDynamicFilters(db, filters)

    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }

	if page <= 0 {
        page = 1
    }

	db = utils.ApplyDynamicSort(db, sorts, "created_at DESC")

	if limit > 0 {
        offset := (page - 1) * limit
        db = db.Offset(offset).Limit(limit)
    }

	if err := db.Find(&locations).Error; err != nil {
        return nil, 0, err
    }

    return locations, int(total), nil
}

func (r *LocationRepository) FindById(locationID uuid.UUID) (*model.Location, error) {
	var location model.Location

	err := r.db.
		First(&location, locationID).
		Error
		
	return &location, err
}

func (r *LocationRepository) Truncate() error {
    var location model.Location

	err := r.db.
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Unscoped().
		Delete(&location).
		Error
    
	return err
}

func (r *LocationRepository) FindByUserID(userID uuid.UUID) (*model.Location, error) {
	var location model.Location

	err := r.db.
		Where("user_id = ?", userID).
		First(&location).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound){
		return nil, err
	}

	return &location, err
}

func (r *LocationRepository) FindAllByUserID(userID uuid.UUID) ([]model.Location, error) {
	var locations []model.Location

	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&locations).
		Error

	return locations, err
}