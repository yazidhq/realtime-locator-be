package repository

import (
	"TeamTrackerBE/internal/domain/model"
	"TeamTrackerBE/internal/utils"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) *ContactRepository {
	return &ContactRepository{db: db}
}

func (r *ContactRepository) Create(contact *model.Contact) (*model.Contact, error) {
	err := r.db.
		Create(contact).
		Error

	return contact, err
}

func (r *ContactRepository) Update(contactID uuid.UUID, req model.Contact) (*model.Contact, error) {
    var contact model.Contact

    if err := r.db.
		First(&contact, contactID).
		Error; err != nil {
			return nil, err
		}

    if err := r.db.
		Model(&contact).
		Updates(req).
		Error; err != nil {
			return nil, err
		}

    return &contact, nil
}

func (r *ContactRepository) Delete(contactID uuid.UUID) (*model.Contact, error) {
	var contact model.Contact

	err := r.db.
		Where("id = ?", contactID).
		Delete(&contact).
		Error

	return &contact, err
}

func (r *ContactRepository) FindAll(page, limit int, filters []utils.FilterOptions) ([]model.Contact, int, error) {
    var contacts []model.Contact
    var total int64

    db := r.db.
		Model(&model.Contact{})

    db = utils.ApplyDynamicFilters(db, filters)

    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    offset := (page - 1) * limit

    if err := db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&contacts).Error; err != nil {
        return nil, 0, err
    }

    return contacts, int(total), nil
}

func (r *ContactRepository) FindById(contactID uuid.UUID) (*model.Contact, error) {
	var contact model.Contact

	err := r.db.
		First(&contact, contactID).
		Error
		
	return &contact, err
}

func (r *ContactRepository) Truncate() error {
    var contact model.Contact

	err := r.db.
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Unscoped().
		Delete(&contact).
		Error
    
	return err
}

func (r *ContactRepository) FindByUserID(userID uuid.UUID) (*model.Contact, error) {
	var contact model.Contact

	err := r.db.
		Where("user_id = ?", userID).
		First(&contact).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound){
		return nil, err
	}

	return &contact, err
}

func (r *ContactRepository) FindByContactID(contactID uuid.UUID) (*model.Contact, error) {
	var contact model.Contact

	err := r.db.
		Where("contact_id = ?", contactID).
		First(&contact).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound){
		return nil, err
	}

	return &contact, err
}

func (r *ContactRepository) FindByUserIDContactID(userID uuid.UUID, contactID uuid.UUID) (*model.Contact, error) {
	var contact model.Contact

	err := r.db.
		Where("user_id = ?", userID).
		Where("contact_id = ?", contactID).
		First(&contact).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound){
		return nil, err
	}

	return &contact, err
}