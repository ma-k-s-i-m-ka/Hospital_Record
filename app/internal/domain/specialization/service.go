package specialization

import (
	"HospitalRecord/app/internal/domain/apperror"
	"HospitalRecord/app/pkg/logger"
	"context"
	"errors"
)

type Service interface {
	Create(ctx context.Context, input *CreateSpecializationDTO) (*Specialization, error)
	GetById(ctx context.Context, id int64) (*Specialization, error)
	Update(ctx context.Context, specialization *UpdateSpecializationDTO) error
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

func (s *service) Create(ctx context.Context, input *CreateSpecializationDTO) (*Specialization, error) {
	s.logger.Info("SERVICE: CREATE SPECIALIZATION")
	specializ := Specialization{
		Name: input.Name,
	}
	specialization, err := s.storage.Create(&specializ)
	if err != nil {
		return nil, err
	}

	return specialization, nil
}

func (s *service) GetById(ctx context.Context, id int64) (*Specialization, error) {
	s.logger.Info("SERVICE: GET SPECIALIZATION BY ID")
	specialization, err := s.storage.FindById(id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.logger.Warnf("cannot find specialization by id: %v", err)
		return nil, err
	}
	return specialization, nil
}

func (s *service) Update(ctx context.Context, specialization *UpdateSpecializationDTO) error {
	s.logger.Info("SERVICE: UPDATE SPECIALIZATION")
	_, err := s.storage.FindById(specialization.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Errorf("failed to get specialization: %v", err)
		}
		return err
	}
	//TODO
	/*if !u.ComparePassword(user.Password) {
		return apperror.ErrWrongPassword
	}*/

	err = s.storage.Update(specialization)
	if err != nil {
		s.logger.Errorf("failed to update user: %v", err)
		return err
	}
	return nil
}

func (s *service) Delete(id int64) error {
	s.logger.Info("SERVICE: DELETE SPECIALIZATION")
	err := s.storage.Delete(id)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Warnf("failed to delete specialization: %v", err)
		}
		return err
	}
	return nil
}
