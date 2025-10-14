package repository

import (
	"TeamTrackerBE/internal/domain/model"
	"TeamTrackerBE/internal/utils"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupParticipantRepository struct {
	db *gorm.DB
}

func NewGroupParticipantRepository(db *gorm.DB) *GroupParticipantRepository {
	return &GroupParticipantRepository{db: db}
}

func (r *GroupParticipantRepository) Create(groupParticipant *model.GroupParticipant) (*model.GroupParticipant, error) {
	err := r.db.
		Create(groupParticipant).
		Error

	return groupParticipant, err
}

func (r *GroupParticipantRepository) Update(groupParticipantID uuid.UUID, req model.GroupParticipant) (*model.GroupParticipant, error) {
    var groupParticipant model.GroupParticipant

    if err := r.db.
		First(&groupParticipant, groupParticipantID).
		Error; err != nil {
			return nil, err
		}

    if err := r.db.
		Model(&groupParticipant).
		Updates(req).
		Error; err != nil {
			return nil, err
		}

    return &groupParticipant, nil
}

func (r *GroupParticipantRepository) Delete(groupParticipantID uuid.UUID) (*model.GroupParticipant, error) {
	var groupParticipant model.GroupParticipant

	err := r.db.
		Where("id = ?", groupParticipantID).
		Delete(&groupParticipant).
		Error

	return &groupParticipant, err
}

func (r *GroupParticipantRepository) FindAll(page, limit int, filters []utils.FilterOptions) ([]model.GroupParticipant, int, error) {
    var groupParticipants []model.GroupParticipant
    var total int64

    db := r.db.
		Model(&model.GroupParticipant{})

    db = utils.ApplyDynamicFilters(db, filters)

    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    offset := (page - 1) * limit

    if err := db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&groupParticipants).Error; err != nil {
        return nil, 0, err
    }

    return groupParticipants, int(total), nil
}

func (r *GroupParticipantRepository) FindById(groupParticipantID uuid.UUID) (*model.GroupParticipant, error) {
	var groupParticipant model.GroupParticipant

	err := r.db.
		First(&groupParticipant, groupParticipantID).
		Error
		
	return &groupParticipant, err
}

func (r *GroupParticipantRepository) Truncate() error {
    var groupParticipant model.GroupParticipant

	err := r.db.
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Unscoped().
		Delete(&groupParticipant).
		Error
    
	return err
}

func (r *GroupParticipantRepository) FindByGroupID(groupID uuid.UUID) (*model.GroupParticipant, error) {
	var groupParticipant model.GroupParticipant

	err := r.db.
		Where("group_id = ?", groupID).
		First(&groupParticipant).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound){
		return nil, err
	}

	return &groupParticipant, err
}

func (r *GroupParticipantRepository) FindByUserID(userID uuid.UUID) (*model.GroupParticipant, error) {
	var groupParticipant model.GroupParticipant

	err := r.db.
		Where("user_id = ?", userID).
		First(&groupParticipant).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound){
		return nil, err
	}

	return &groupParticipant, err
}