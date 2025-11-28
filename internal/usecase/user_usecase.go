package usecase

import (
	"TeamTrackerBE/internal/delivery/http/dto"
	"TeamTrackerBE/internal/domain/model"
	"TeamTrackerBE/internal/domain/repository"
	"TeamTrackerBE/internal/utils"
	"TeamTrackerBE/internal/utils/responses"

	"github.com/google/uuid"
)

type UserUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(r *repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: *r}
}

func (u *UserUsecase) Create(req *dto.UserCreateRequest) (*model.User, error) {
	user := &model.User{
		Role: model.Role(req.Role),
		Name: req.Name,
		Username: req.Username,
		Email: req.Email,
		PhoneNumber: req.PhoneNumber,
		Password: req.Password,
	}

	existingEmail, _ := u.repo.FindByEmail(user.Email)
	if existingEmail != nil {
		return nil, responses.NewBadRequestError("email already exist")
	}
	
	existingUsername, _ := u.repo.FindByUsername(user.Username)
	if existingUsername != nil {
		return nil, responses.NewBadRequestError("username already exist")
	}
	
	existingPhoneNumber, _ := u.repo.FindByPhoneNumber(user.PhoneNumber)
	if existingPhoneNumber != nil {
		return nil, responses.NewBadRequestError("phone number already exist")
	}

	hashedPassword, errHash := utils.HashPassword(user.Password)
	if errHash != nil {
		return nil, errHash
	}

	user.Password = hashedPassword
	created, err := u.repo.Create(user)
	if err != nil {
		return nil, err
	}
	
	return created, nil
}

func (u *UserUsecase) Update(userID uuid.UUID, req *dto.UserUpdateRequest) (*model.User, error) {
	user := &model.User{
		Role: model.Role(req.Role),
		Name: req.Name,
		Username: req.Username,
		Email: req.Email,
		PhoneNumber: req.PhoneNumber,
		Password: req.Password,
	}

	if _, err := u.repo.FindById(userID); err != nil {
		return nil, responses.NewNotFoundError("user not found")
	}

	if user.Email != "" {
		existingEmail, _ := u.repo.FindByEmail(user.Email)
		if existingEmail != nil && existingEmail.ID != userID {
			return nil, responses.NewBadRequestError("email already exists")
		}
	}
	
	if user.Username != "" {
		existingUsername, _ := u.repo.FindByUsername(user.Username)
		if existingUsername != nil && existingUsername.ID != userID {
			return nil, responses.NewBadRequestError("username already exists")
		}
	}
	
	if user.PhoneNumber != "" {
		existingPhoneNumber, _ := u.repo.FindByPhoneNumber(user.PhoneNumber)
		if existingPhoneNumber != nil && existingPhoneNumber.ID != userID {
			return nil, responses.NewBadRequestError("phone number already exists")
		}
	}
	
	if user.Password != "" {
		hashedPassword, errHash := utils.HashPassword(user.Password)
		if errHash != nil {
			return nil, errHash
		}
		user.Password = hashedPassword
	}

	updated, err := u.repo.Update(userID, *user)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (u *UserUsecase) Delete(userID uuid.UUID) (*model.User, error) {
	if _, err := u.repo.FindById(userID); err != nil {
		return nil, responses.NewNotFoundError("user not found")
	}
	
	deleted, err := u.repo.Delete(userID)
	if err != nil {
		return nil, err
	}

	return deleted, nil
}

func (u *UserUsecase) FindAll(page, limit int, filters []utils.FilterOptions, sorts []utils.SortOption) ([]model.User, int, error) {
	result, total, err := u.repo.FindAll(page, limit, filters, sorts)
    if err != nil {
        return nil, 0, err
    }
    return result, total, nil
}

func (u *UserUsecase) FindById(userID uuid.UUID) (*model.User, error) {
	result, err := u.repo.FindById(userID)
	if err != nil {
		return nil, responses.NewNotFoundError("user not found")
	}

	return result, nil
}

func (u *UserUsecase) Truncate() error {
	err := u.repo.Truncate()
	if err != nil {
		return err
	}
	return nil
}
