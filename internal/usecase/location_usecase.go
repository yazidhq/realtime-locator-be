package usecase

import (
	"TeamTrackerBE/internal/delivery/http/dto"
	"TeamTrackerBE/internal/domain/model"
	"TeamTrackerBE/internal/domain/repository"
	"TeamTrackerBE/internal/utils"
	"TeamTrackerBE/internal/utils/responses"
	"time"

	"github.com/google/uuid"
)

type LocationUsecase struct {
	repo repository.LocationRepository
	repoUser repository.UserRepository
}

func NewLocationUsecase(
	r *repository.LocationRepository,
	ru *repository.UserRepository,
	) *LocationUsecase {
	return &LocationUsecase{
		repo: *r,
		repoUser: *ru,
	}
}

func (u *LocationUsecase) Create(req *dto.LocationCreateRequest) (*model.Location, error) {
	location := &model.Location{
		UserID: req.UserID,
		Latitude: req.Latitude,
		Longitude: req.Longitude,
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
		UserID: req.UserID,
		Latitude: req.Latitude,
		Longitude: req.Longitude,
	}

	if _, err := u.repo.FindById(locationID); err != nil {
		return nil, responses.NewNotFoundError("location not found")
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

func (u *LocationUsecase) FindAll(page, limit int, filters []utils.FilterOptions, sorts []utils.SortOption) ([]model.Location, int, error) {
	result, total, err := u.repo.FindAll(page, limit, filters, sorts)
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

func (u *LocationUsecase) HistoryByUser(userID uuid.UUID) ([]dto.LocationHistoryGroupResponse, error) {
	if _, err := u.repoUser.FindById(userID); err != nil {
		return nil, responses.NewBadRequestError("user id not found in user")
	}

	locations, err := u.repo.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	jakartaLoc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		jakartaLoc = time.Local
	}

	todayKey := time.Now().In(jakartaLoc).Format("02-01-2006")

	groups := make([]dto.LocationHistoryGroupResponse, 0)
	var currentDate string
	var currentIndex int = -1

	for _, loc := range locations {
		locTime := loc.CreatedAt.In(jakartaLoc)
		dateKey := locTime.Format("02-01-2006")
		if dateKey == todayKey {
			continue
		}
		if dateKey != currentDate {
			currentDate = dateKey
			groups = append(groups, dto.LocationHistoryGroupResponse{
				Date:      dateKey,
				Locations: []dto.LocationHistoryItemResponse{},
			})
			currentIndex++
		}

		groups[currentIndex].Locations = append(groups[currentIndex].Locations, dto.LocationHistoryItemResponse{
			ID:        loc.ID.String(),
			UserID:    loc.UserID,
			Latitude:  loc.Latitude,
			Longitude: loc.Longitude,
			CreatedAt: locTime,
		})
	}

	return groups, nil
}
