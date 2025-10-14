package usecase

import (
	"TeamTrackerBE/internal/delivery/http/dto"
	"TeamTrackerBE/internal/domain/model"
	"TeamTrackerBE/internal/domain/repository"
	"TeamTrackerBE/internal/utils"
	"TeamTrackerBE/internal/utils/responses"

	"github.com/google/uuid"
)

type ContactUsecase struct {
	repo repository.ContactRepository
	repoUser repository.UserRepository
}

func NewContactUsecase(r *repository.ContactRepository) *ContactUsecase {
	return &ContactUsecase{repo: *r}
}

func (u *ContactUsecase) Create(req *dto.ContactCreateRequest) (*model.Contact, error) {
	contact := &model.Contact{
		UserID: req.UserID,
		ContactID: req.ContactID,
		Status: req.Status,
	}

	if _, err := u.repoUser.FindById(contact.UserID); err != nil {
		return nil, responses.NewBadRequestError("user id not found in user")
	}
	
	if _, err := u.repoUser.FindById(contact.ContactID); err != nil {
		return nil, responses.NewBadRequestError("contact id not found in user")
	}

	created, err := u.repo.Create(contact)
	if err != nil {
		return nil, err
	}
	
	return created, nil
}

func (u *ContactUsecase) Update(contactID uuid.UUID, req *dto.ContactUpdateRequest) (*model.Contact, error) {
	contact := &model.Contact{
		UserID: req.UserID,
		ContactID: req.ContactID,
		Status: req.Status,
	}

	if _, err := u.repo.FindById(contactID); err != nil {
		return nil, responses.NewNotFoundError("contact not found")
	}

	if contact.UserID.String() != "" {
		if _, err := u.repoUser.FindById(contact.UserID); err != nil {
			return nil, responses.NewBadRequestError("user id not found in user")
		}
	}
	
	if contact.ContactID.String() != "" {
		if _, err := u.repoUser.FindById(contact.ContactID); err != nil {
			return nil, responses.NewBadRequestError("contact id not found in user")
		}
	}

	updated, err := u.repo.Update(contactID, *contact)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (u *ContactUsecase) Delete(contactID uuid.UUID) (*model.Contact, error) {
	if _, err := u.repo.FindById(contactID); err != nil {
		return nil, responses.NewNotFoundError("contact not found")
	}
	
	deleted, err := u.repo.Delete(contactID)
	if err != nil {
		return nil, err
	}

	return deleted, nil
}

func (u *ContactUsecase) FindAll(page, limit int, filters []utils.FilterOptions) ([]model.Contact, int, error) {
    result, total, err := u.repo.FindAll(page, limit, filters)
    if err != nil {
        return nil, 0, err
    }
    return result, total, nil
}

func (u *ContactUsecase) FindById(contactID uuid.UUID) (*model.Contact, error) {
	result, err := u.repo.FindById(contactID)
	if err != nil {
		return nil, responses.NewNotFoundError("contact not found")
	}

	return result, nil
}

func (u *ContactUsecase) Truncate() error {
	err := u.repo.Truncate()
	if err != nil {
		return err
	}
	return nil
}
