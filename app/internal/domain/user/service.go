package user

import (
	"HospitalRecord/app/internal/domain/apperror"
	"HospitalRecord/app/pkg/logger"
	"context"
	"errors"
	"fmt"
)

type Service interface {
	Create(ctx context.Context, user *CreateUserDTO) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetById(ctx context.Context, id int64) (*User, error)
	GetByPolicyNumber(ctx context.Context, policy string) (*User, error)
	Update(ctx context.Context, user *UpdateUserDTO) error
	PartiallyUpdate(ctx context.Context, user *PartiallyUpdateUserDTO) error
	Delete(id int64) error
}

type service struct {
	logger  logger.Logger
	storage Storage
}

func NewService(storage Storage, logger logger.Logger) Service {
	return &service{
		logger:  logger,
		storage: storage,
	}
}
func (s *service) Create(ctx context.Context, input *CreateUserDTO) (*User, error) {
	s.logger.Info("SERVICE: CREATE USER")
	checkEmail, err := s.storage.FindByEmail(input.Email)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return nil, err
		}
	}

	if checkEmail != nil {
		return nil, apperror.ErrRepeatedEmail
	}

	u := User{
		Email:        input.Email,
		Name:         input.Name,
		Surname:      input.Surname,
		Age:          input.Age,
		Gender:       input.Gender,
		Password:     input.Password,
		PolicyNumber: input.PolicyNumber,
	}

	err = u.HashPassword()
	if err != nil {
		return nil, fmt.Errorf("cannot hash password")
	}

	user, err := s.storage.Create(&u)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) GetByEmail(ctx context.Context, email string) (*User, error) {
	s.logger.Info("SERVICE: GET USER BY EMAIL")
	user, err := s.storage.FindByEmail(email)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.logger.Warnf("cannot find user by email: %v", err)
		return nil, err
	}
	return user, nil
}

func (s *service) GetById(ctx context.Context, id int64) (*User, error) {
	s.logger.Info("SERVICE: GET USER BY ID")
	user, err := s.storage.FindById(id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.logger.Warnf("cannot find user by id: %v", err)
		return nil, err
	}
	return user, nil
}

func (s *service) GetByPolicyNumber(ctx context.Context, policy string) (*User, error) {
	s.logger.Info("SERVICE: GET USER BY POLICY NUMBER")
	user, err := s.storage.FindByPolicyNumber(policy)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.logger.Warnf("cannot find user by policy number: %v", err)
		return nil, err
	}
	return user, nil
}

func (s *service) Update(ctx context.Context, user *UpdateUserDTO) error {
	s.logger.Info("SERVICE: UPDATE USER")
	u, err := s.storage.FindById(user.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Errorf("failed to get user: %v", err)
		}
		return err
	}
	u.Password = user.Password
	err = user.HashPassword()
	if err != nil {
		s.logger.Errorf("failed to hash password: %v", err)
		return err
	}

	err = s.storage.Update(user)
	if err != nil {
		s.logger.Errorf("failed to update user: %v", err)
		return err
	}
	return nil
}

func (s *service) PartiallyUpdate(ctx context.Context, user *PartiallyUpdateUserDTO) error {
	s.logger.Info("SERVICE: PARTIALLY UPDATE USER")
	u, err := s.storage.FindById(user.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Errorf("failed to get user: %v", err)
		}
		return err
	}

	if user.Password != nil {
		u.Password = *user.Password
		err = user.HashPassword()
		if err != nil {
			s.logger.Errorf("failed to hash password: %v", err)
			return err
		}
	}

	err = s.storage.PartiallyUpdate(user)
	if err != nil {
		s.logger.Errorf("failed to partially update user: %v", err)
		return err
	}
	return nil
}

func (s *service) Delete(id int64) error {
	s.logger.Info("SERVICE: DELETE USER")
	err := s.storage.Delete(id)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Warnf("failed to delete user: %v", err)
		}
		return err
	}
	return nil
}
