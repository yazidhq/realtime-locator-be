package usecase

import (
	"TeamTrackerBE/internal/delivery/http/dto"
	"TeamTrackerBE/internal/domain/model"
	"TeamTrackerBE/internal/domain/repository"
	"TeamTrackerBE/internal/utils"
	"TeamTrackerBE/internal/utils/responses"

	"github.com/google/uuid"
)

type LocationUsecase struct {
	repo repository.LocationRepository
	repoGroup repository.GroupRepository
	repoUser repository.UserRepository
}

func NewLocationUsecase(r *repository.LocationRepository) *LocationUsecase {
	return &LocationUsecase{repo: *r}
}

func (u *LocationUsecase) Create(req *dto.LocationCreateRequest) (*model.Location, error) {
	location := &model.Location{
		GroupID: req.GroupID,
		UserID: req.UserID,
		Latitude: req.Latitude,
		Longitude: req.Longitude,
	}

	if _, err := u.repoGroup.FindById(location.GroupID); err != nil {
		return nil, responses.NewBadRequestError("group id not found in group")
	}

	if _, err := u.repoUser.FindById(location.UserID); err != nil {
		return nil, responses.NewBadRequestError("user id not found in user")
	}

	created, err := u.repo.Create(location)
	if err != nil {
		return nil, err
	}
	
	return created, nil
}

func (u *LocationUsecase) Update(locationID uuid.UUID, req *dto.LocationUpdateRequest) (*model.Location, error) {
	location := &model.Location{
		GroupID: req.GroupID,
		UserID: req.UserID,
		Latitude: req.Latitude,
		Longitude: req.Longitude,
	}

	if _, err := u.repo.FindById(locationID); err != nil {
		return nil, responses.NewNotFoundError("location not found")
	}

	if location.GroupID.String() != "" {
		if _, err := u.repoGroup.FindById(location.GroupID); err != nil {
			return nil, responses.NewBadRequestError("group id not found in group")
		}
	}

	if location.UserID.String() != "" {
		if _, err := u.repoUser.FindById(location.UserID); err != nil {
			return nil, responses.NewBadRequestError("user id not found in user")
		}
	}

	updated, err := u.repo.Update(locationID, *location)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (u *LocationUsecase) Delete(locationID uuid.UUID) (*model.Location, error) {
	if _, err := u.repo.FindById(locationID); err != nil {
		return nil, responses.NewNotFoundError("location not found")
	}
	
	deleted, err := u.repo.Delete(locationID)
	if err != nil {
		return nil, err
	}

	return deleted, nil
}

func (u *LocationUsecase) FindAll(page, limit int, filters []utils.FilterOptions) ([]model.Location, int, error) {
    result, total, err := u.repo.FindAll(page, limit, filters)
    if err != nil {
        return nil, 0, err
    }
    return result, total, nil
}

func (u *LocationUsecase) FindById(locationID uuid.UUID) (*model.Location, error) {
	result, err := u.repo.FindById(locationID)
	if err != nil {
		return nil, responses.NewNotFoundError("location not found")
	}

	return result, nil
}

func (u *LocationUsecase) Truncate() error {
	err := u.repo.Truncate()
	if err != nil {
		return err
	}
	return nil
}
