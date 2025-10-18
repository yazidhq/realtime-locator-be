package usecase

import (
	"TeamTrackerBE/internal/delivery/http/dto"
	"TeamTrackerBE/internal/domain/model"
	"TeamTrackerBE/internal/domain/repository"
	"TeamTrackerBE/internal/utils"
	"TeamTrackerBE/internal/utils/responses"
	"encoding/json"

	"github.com/google/uuid"
)

type GroupUsecase struct {
	repo repository.GroupRepository
	repoUser repository.UserRepository
}

func NewGroupUsecase(
	r *repository.GroupRepository,
	ru *repository.UserRepository,
	) *GroupUsecase {
	return &GroupUsecase{
		repo: *r,
		repoUser: *ru,
	}
}

func (u *GroupUsecase) Create(req *dto.GroupCreateRequest) (*model.Group, error) {
	group := &model.Group{
		Name: req.Name,
		OwnerID: req.OwnerID,
		RadiusArea: req.RadiusArea,
	}

	if len(group.RadiusArea) > 0 {
		var area dto.RadiusArea
		if err := json.Unmarshal(group.RadiusArea, &area); err != nil {
			return nil, responses.NewBadRequestError("radius area harus berupa JSON dengan format yang benar")
		}

		if area.Radius <= 0 {
			return nil, responses.NewBadRequestError("radius harus lebih besar dari 0")
		}

		if area.CenterLat == 0 || area.CenterLon == 0 {
			return nil, responses.NewBadRequestError("area harus memiliki center lat dan center lon yang valid")
		}
	}

	if _, err := u.repoUser.FindById(group.OwnerID); err != nil {
		return nil, responses.NewBadRequestError("owner id not found in user")
	}

	created, err := u.repo.Create(group)
	if err != nil {
		return nil, err
	}
	
	return created, nil
}

func (u *GroupUsecase) Update(groupID uuid.UUID, req *dto.GroupUpdateRequest) (*model.Group, error) {
	group := &model.Group{
		Name: req.Name,
		OwnerID: req.OwnerID,
		RadiusArea: req.RadiusArea,
	}

	if _, err := u.repo.FindById(groupID); err != nil {
		return nil, responses.NewNotFoundError("group not found")
	}

	if group.OwnerID.String() != "" {
		if _, err := u.repoUser.FindById(group.OwnerID); err != nil {
			return nil, responses.NewBadRequestError("owner id not found in user")
		}
	}

	updated, err := u.repo.Update(groupID, *group)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (u *GroupUsecase) Delete(groupID uuid.UUID) (*model.Group, error) {
	if _, err := u.repo.FindById(groupID); err != nil {
		return nil, responses.NewNotFoundError("group not found")
	}
	
	deleted, err := u.repo.Delete(groupID)
	if err != nil {
		return nil, err
	}

	return deleted, nil
}

func (u *GroupUsecase) FindAll(page, limit int, filters []utils.FilterOptions) ([]model.Group, int, error) {
    result, total, err := u.repo.FindAll(page, limit, filters)
    if err != nil {
        return nil, 0, err
    }
    return result, total, nil
}

func (u *GroupUsecase) FindById(groupID uuid.UUID) (*model.Group, error) {
	result, err := u.repo.FindById(groupID)
	if err != nil {
		return nil, responses.NewNotFoundError("group not found")
	}

	return result, nil
}

func (u *GroupUsecase) Truncate() error {
	err := u.repo.Truncate()
	if err != nil {
		return err
	}
	return nil
}
