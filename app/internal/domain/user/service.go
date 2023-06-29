package user

import (
	"HospitalRecord/app/internal/domain/apperror"
	"HospitalRecord/app/pkg/logger"
	"context"
	"errors"
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
func (s *service) Create(ctx context.Context, v *CreateUserDTO) (*User, error) {
	checkEmail, err := s.storage.FindByEmail(v.Email)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return nil, err
		}
	}

	if checkEmail != nil {
		return nil, apperror.ErrRepeatedEmail
	}

	u := User{
		Email:        v.Email,
		Name:         v.Name,
		Surname:      v.Surname,
		Age:          v.Age,
		Gender:       v.Gender,
		Password:     v.Password,
		PolicyNumber: v.PolicyNumber,
	}
	//TODO
	/*	err = u.HashPassword()
		if err != nil {
			return nil, fmt.Errorf("cannot hash password")
		}
	*/
	user, err := s.storage.Create(&u)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) GetByEmail(ctx context.Context, email string) (*User, error) {
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
	_, err := s.GetById(ctx, user.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Errorf("failed to get user: %v", err)
		}
		return err
	}
	//TODO
	/*if !u.ComparePassword(user.Password) {
		return apperror.ErrWrongPassword
	}*/

	err = s.storage.Update(user)
	if err != nil {
		s.logger.Errorf("failed to update user: %v", err)
		return err
	}
	return nil
}

func (s *service) PartiallyUpdate(ctx context.Context, user *PartiallyUpdateUserDTO) error {
	_, err := s.GetById(ctx, user.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Errorf("failed to get user: %v", err)
		}
		return err
	}

	/*if !u.ComparePassword(*user.OldPassword) {
		return apperror.ErrWrongPassword
	}

	if user.NewPassword != nil {
		u.Password = *user.NewPassword
		err = user.HashPassword()
		if err != nil {
			s.logger.Errorf("failed ot hash password: %v", err)
			return err
		}
	}
	*/
	err = s.storage.PartiallyUpdate(user)
	if err != nil {
		s.logger.Errorf("failed to partially update user: %v", err)
		return err
	}
	return nil
}

func (s *service) Delete(id int64) error {
	err := s.storage.Delete(id)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Warnf("failed to delete user: %v", err)
		}
		return err
	}
	return nil
}
