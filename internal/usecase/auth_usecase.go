package usecase

import (
	"TeamTrackerBE/internal/delivery/http/dto"
	"TeamTrackerBE/internal/domain/model"
	"TeamTrackerBE/internal/domain/repository"
	"TeamTrackerBE/internal/utils"
	"errors"
	"strings"
)

type AuthInterface interface {
	Register(req *dto.RegisterRequest) (*model.User, error)
	Login(email, password string) (string, *model.User, error)
}

type AuthUsecase struct {
	repo repository.UserRepository
}

func NewAuthUsecase(r *repository.UserRepository) *AuthUsecase {
	return &AuthUsecase{repo: *r}
}

func (u *AuthUsecase) Register(req *dto.RegisterRequest) (*model.User, error) {
	existingEmail, _ := u.repo.FindByEmail(req.Email)
	if existingEmail != nil {
		return nil, errors.New("email already exists")
	}
	
	existingUsername, _ := u.repo.FindByUsername(req.Username)
	if existingUsername != nil {
		return nil, errors.New("username already exists")
	}
	
	existingPhoneNumber, _ := u.repo.FindByPhoneNumber(req.PhoneNumber)
	if existingPhoneNumber != nil {
		return nil, errors.New("phone number already exists")
	}

	hashedPassword, errHash := utils.HashPassword(req.Password)
	if errHash != nil {
		return nil, errHash
	}

	user := &model.User{
		Name: req.Name,
		Username: req.Username,
		Email: req.Email,
		PhoneNumber: req.PhoneNumber,
		Password: hashedPassword,
	}

	created, errCreate := u.repo.Create(user)
	if errCreate != nil {
		return nil, errCreate
	}

	return created, nil
}

func (u *AuthUsecase) Login(email, password string) (string, *model.User, error) {
	email = strings.ToLower(email)

	found, err := u.repo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	if !utils.VerifyPassword(password, found.Password) {
		return "", nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateJWT(found.ID, model.Role(found.Role))
	if err != nil {
		return "", nil, err
	}

	return token, found, nil
}