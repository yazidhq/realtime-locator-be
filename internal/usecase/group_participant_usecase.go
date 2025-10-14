package usecase

import (
	"TeamTrackerBE/internal/delivery/http/dto"
	"TeamTrackerBE/internal/domain/model"
	"TeamTrackerBE/internal/domain/repository"
	"TeamTrackerBE/internal/utils"
	"TeamTrackerBE/internal/utils/responses"

	"github.com/google/uuid"
)

type GroupParticipantUsecase struct {
	repo repository.GroupParticipantRepository
	repoGroup repository.GroupRepository
	repoUser repository.UserRepository
}

func NewGroupParticipantUsecase(r *repository.GroupParticipantRepository) *GroupParticipantUsecase {
	return &GroupParticipantUsecase{repo: *r}
}

func (u *GroupParticipantUsecase) Create(req *dto.GroupParticipantCreateRequest) (*model.GroupParticipant, error) {
	groupParticipant := &model.GroupParticipant{
		GroupID: req.GroupID,
		UserID: req.UserID,
	}

	if _, err := u.repoGroup.FindById(groupParticipant.GroupID); err != nil {
		return nil, responses.NewBadRequestError("group id not found in group")
	}

	if _, err := u.repoUser.FindById(groupParticipant.UserID); err != nil {
		return nil, responses.NewBadRequestError("user id not found in user")
	}

	created, err := u.repo.Create(groupParticipant)
	if err != nil {
		return nil, err
	}
	
	return created, nil
}

func (u *GroupParticipantUsecase) Update(groupParticipantID uuid.UUID, req *dto.GroupParticipantUpdateRequest) (*model.GroupParticipant, error) {
	groupParticipant := &model.GroupParticipant{
		GroupID: req.GroupID,
		UserID: req.UserID,
	}

	if _, err := u.repo.FindById(groupParticipantID); err != nil {
		return nil, responses.NewNotFoundError("groupParticipant not found")
	}

	if groupParticipant.GroupID.String() != "" {
		if _, err := u.repoGroup.FindById(groupParticipant.GroupID); err != nil {
			return nil, responses.NewBadRequestError("group id not found in group")
		}
	}

	if groupParticipant.UserID.String() != "" {
		if _, err := u.repoUser.FindById(groupParticipant.UserID); err != nil {
			return nil, responses.NewBadRequestError("user id not found in user")
		}
	}

	updated, err := u.repo.Update(groupParticipantID, *groupParticipant)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (u *GroupParticipantUsecase) Delete(groupParticipantID uuid.UUID) (*model.GroupParticipant, error) {
	if _, err := u.repo.FindById(groupParticipantID); err != nil {
		return nil, responses.NewNotFoundError("group participant not found")
	}
	
	deleted, err := u.repo.Delete(groupParticipantID)
	if err != nil {
		return nil, err
	}

	return deleted, nil
}

func (u *GroupParticipantUsecase) FindAll(page, limit int, filters []utils.FilterOptions) ([]model.GroupParticipant, int, error) {
    result, total, err := u.repo.FindAll(page, limit, filters)
    if err != nil {
        return nil, 0, err
    }
    return result, total, nil
}

func (u *GroupParticipantUsecase) FindById(groupParticipantID uuid.UUID) (*model.GroupParticipant, error) {
	result, err := u.repo.FindById(groupParticipantID)
	if err != nil {
		return nil, responses.NewNotFoundError("group participant not found")
	}

	return result, nil
}

func (u *GroupParticipantUsecase) Truncate() error {
	err := u.repo.Truncate()
	if err != nil {
		return err
	}
	return nil
}
