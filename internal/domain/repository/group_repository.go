package repository

import (
	"TeamTrackerBE/internal/domain/model"
	"TeamTrackerBE/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(group *model.Group) (*model.Group, error) {
	err := r.db.
		Create(group).
		Error

	return group, err
}

func (r *GroupRepository) Update(groupID uuid.UUID, req model.Group) (*model.Group, error) {
    var group model.Group

    if err := r.db.
		First(&group, groupID).
		Error; err != nil {
			return nil, err
		}

    if err := r.db.
		Model(&group).
		Updates(req).
		Error; err != nil {
			return nil, err
		}

    return &group, nil
}

func (r *GroupRepository) Delete(groupID uuid.UUID) (*model.Group, error) {
	var group model.Group

	err := r.db.
		Where("id = ?", groupID).
		Delete(&group).
		Error

	return &group, err
}

func (r *GroupRepository) FindAll(page, limit int, filters []utils.FilterOptions) ([]model.Group, int, error) {
    var groups []model.Group
    var total int64

    db := r.db.
		Model(&model.Group{})

    db = utils.ApplyDynamicFilters(db, filters)

    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    offset := (page - 1) * limit

    if err := db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&groups).Error; err != nil {
        return nil, 0, err
    }

    return groups, int(total), nil
}

func (r *GroupRepository) FindById(groupID uuid.UUID) (*model.Group, error) {
	var group model.Group

	err := r.db.
		First(&group, groupID).
		Error
		
	return &group, err
}

func (r *GroupRepository) Truncate() error {
    var group model.Group

	err := r.db.
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Unscoped().
		Delete(&group).
		Error
    
	return err
}