package doctor

import (
	"HospitalRecord/app/internal/domain/apperror"
	"HospitalRecord/app/pkg/logger"
	"context"
	"errors"
)

type Service interface {
	Create(ctx context.Context, doctor *CreateDoctorDTO) (*Doctor, error)
	FindAll() (*[]Doctor, error)
	FindAllAvailable(ctx context.Context, id int64, recordingIsAvailable bool) (*[]Doctor, error)
	GetById(ctx context.Context, id int64) (*Doctor, error)
	GetByPortfolioId(ctx context.Context, id int64) (*Doctor, error)
	Update(ctx context.Context, doctor *UpdateDoctorDTO) error
	PartiallyUpdate(ctx context.Context, doctor *PartiallyUpdateDoctorDTO) error
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

func (s *service) Create(ctx context.Context, d *CreateDoctorDTO) (*Doctor, error) {
	checkPortfolioID, err := s.storage.FindByPortfolioId(d.PortfolioID)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return nil, err
		}
	}

	if checkPortfolioID != nil {
		return nil, apperror.ErrRepeatedPortfolioId
	}
	doc := Doctor{
		Name:             d.Name,
		Surname:          d.Surname,
		Patronymic:       d.Patronymic,
		ImageID:          d.ImageID,
		Gender:           d.Gender,
		Rating:           d.Rating,
		Age:              d.Age,
		SpecializationID: d.SpecializationID,
		PortfolioID:      d.PortfolioID,
	}
	doctor, err := s.storage.Create(&doc)
	if err != nil {
		return nil, err
	}
	return doctor, nil
}

func (s *service) FindAll() (*[]Doctor, error) {
	doctor, err := s.storage.FindAll()
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.logger.Warnf("cannot find doctors by id: %v", err)
		return nil, err
	}
	return &doctor, nil
	//TODO
}

func (s *service) FindAllAvailable(ctx context.Context, id int64, recordingIsAvailable bool) (*[]Doctor, error) {
	doctor, err := s.storage.FindAllAvailable(id, recordingIsAvailable)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.logger.Warnf("cannot find available doctor by id: %v", err)
		return nil, err
	}
	return &doctor, nil
	//TODO
}

func (s *service) GetById(ctx context.Context, id int64) (*Doctor, error) {
	doctor, err := s.storage.FindById(id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.logger.Warnf("cannot find doctor by id: %v", err)
		return nil, err
	}
	return doctor, nil
}

func (s *service) GetByPortfolioId(ctx context.Context, id int64) (*Doctor, error) {
	doctor, err := s.storage.FindByPortfolioId(id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.logger.Warnf("cannot find doctor by portfolio id: %v", err)
		return nil, err
	}
	return doctor, nil
}

func (s *service) Update(ctx context.Context, doctor *UpdateDoctorDTO) error {
	_, err := s.GetById(ctx, doctor.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Errorf("failed to get doctor: %v", err)
		}
		return err
	}

	err = s.storage.Update(doctor)
	if err != nil {
		s.logger.Errorf("failed to update doctor: %v", err)
		return err
	}
	return nil
}

func (s *service) PartiallyUpdate(ctx context.Context, doctor *PartiallyUpdateDoctorDTO) error {
	_, err := s.GetById(ctx, doctor.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Errorf("failed to get doctor: %v", err)
		}
		return err
	}

	err = s.storage.PartiallyUpdate(doctor)
	if err != nil {
		s.logger.Errorf("failed to partially update doctor: %v", err)
		return err
	}
	return nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	err := s.storage.Delete(id)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Warnf("failed to delete doctor: %v", err)
		}
		return err
	}
	return nil
}
