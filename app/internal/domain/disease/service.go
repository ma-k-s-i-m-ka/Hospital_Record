package disease

import (
	"HospitalRecord/app/internal/domain/apperror"
	"HospitalRecord/app/pkg/logger"
	"context"
	"errors"
)

type Service interface {
	Create(ctx context.Context, input *CreateDiseaseDTO) (*Disease, error)
	GetById(ctx context.Context, id int64) (*Disease, error)
	Update(ctx context.Context, disease *UpdateDiseaseDTO) error
	Delete(ctx context.Context, id int64) error
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

func (s *service) Create(ctx context.Context, input *CreateDiseaseDTO) (*Disease, error) {
	dis := Disease{
		BodyPart:    input.BodyPart,
		Description: input.Description,
	}
	disease, err := s.storage.Create(&dis)
	if err != nil {
		return nil, err
	}

	return disease, nil
}

func (s *service) GetById(ctx context.Context, id int64) (*Disease, error) {
	disease, err := s.storage.FindById(id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.logger.Warnf("cannot find disease by id: %v", err)
		return nil, err
	}
	return disease, nil
}

func (s *service) Update(ctx context.Context, disease *UpdateDiseaseDTO) error {
	_, err := s.GetById(ctx, disease.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Errorf("failed to get disease: %v", err)
		}
		return err
	}
	err = s.storage.Update(disease)
	if err != nil {
		s.logger.Errorf("failed to update disease: %v", err)
		return err
	}
	return nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	err := s.storage.Delete(id)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Warnf("failed to delete disease: %v", err)
		}
		return err
	}
	return nil
}
