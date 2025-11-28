package repository

import (
	"TeamTrackerBE/internal/domain/model"
	"TeamTrackerBE/internal/utils"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) (*model.User, error) {
	err := r.db.
		Create(user).
		Error

	return user, err
}

func (r *UserRepository) BulkCreate(users []model.User) error {
	if len(users) == 0 {
		return nil
	}
	
	return r.db.CreateInBatches(&users, 200).Error
}

func (r *UserRepository) Update(userID uuid.UUID, req model.User) (*model.User, error) {
    var user model.User

    if err := r.db.
		First(&user, userID).
		Error; err != nil {
			return nil, err
		}

    if err := r.db.
		Model(&user).
		Updates(req).
		Error; err != nil {
			return nil, err
		}

    return &user, nil
}

func (r *UserRepository) Delete(userID uuid.UUID) (*model.User, error) {
	var user model.User

	err := r.db.
		Where("id = ?", userID).
		Unscoped().
		Delete(&user).
		Error

	return &user, err
}

func (r *UserRepository) FindAll(page, limit int, filters []utils.FilterOptions, sorts []utils.SortOption) ([]model.User, int, error) {
    var users []model.User
    var total int64

    db := r.db.
		Model(&model.User{})

    db = utils.ApplyDynamicFilters(db, filters)

    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    offset := (page - 1) * limit
	
	db = utils.ApplyDynamicSort(db, sorts, "created_at DESC")

	if err := db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
        return nil, 0, err
    }

    return users, int(total), nil
}

func (r *UserRepository) FindById(userID uuid.UUID) (*model.User, error) {
	var user model.User

	err := r.db.
		First(&user, userID).
		Error
		
	return &user, err
}

func (r *UserRepository) Truncate() error {
    var user model.User

	err := r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&user).Error
    
	return err
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	err := r.db.
		Where("email = ?", email).
		First(&user).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound){
		return nil, err
	}

	return &user, err
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User

	err := r.db.
		Where("username = ?", username).
		First(&user).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound){
		return nil, err
	}

	return &user, err
}

func (r *UserRepository) FindByPhoneNumber(phoneNumber string) (*model.User, error) {
	var user model.User

	err := r.db.
		Where("phone_number = ?", phoneNumber).
		First(&user).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &user, err
}
