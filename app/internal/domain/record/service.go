package record

import (
	"HospitalRecord/app/internal/domain/apperror"
	"HospitalRecord/app/pkg/logger"
	"context"
	"errors"
)

type Service interface {
	Create(ctx context.Context, record *CreateRecordDTO) (*Record, error)
	GetByPatientsId(ctx context.Context, id int64) (*Record, error)
	GetById(ctx context.Context, id int64) (*Record, error)
	Update(ctx context.Context, record *UpdateRecordDTO) error
	PartiallyUpdate(ctx context.Context, record *PartiallyUpdateRecordDTO) error
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
func (s *service) Create(ctx context.Context, input *CreateRecordDTO) (*Record, error) {
	s.logger.Info("SERVICE: CREATE RECORD")

	r := Record{
		ID:               input.ID,
		HospitalAddress:  input.HospitalAddress,
		DoctorOffice:     input.DoctorOffice,
		Tagging:          input.Tagging,
		PatientsID:       input.PatientsID,
		DoctorID:         input.DoctorID,
		SpecializationID: input.SpecializationID,
		TimeRecord:       input.TimeRecord,
	}

	record, err := s.storage.CreateRecord(&r)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (s *service) GetByPatientsId(ctx context.Context, id int64) (*Record, error) {
	s.logger.Info("SERVICE: GET RECORD BY PATIENTS ID")
	record, err := s.storage.FindRecordByPatientsId(id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.logger.Warnf("cannot find record by patients id: %v", err)
		return nil, err
	}
	return record, nil
}

func (s *service) GetById(ctx context.Context, id int64) (*Record, error) {
	s.logger.Info("SERVICE: GET RECORD BY ID")
	record, err := s.storage.FindRecordById(id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.logger.Warnf("cannot find record by id: %v", err)
		return nil, err
	}
	return record, nil
}

func (s *service) Update(ctx context.Context, record *UpdateRecordDTO) error {
	s.logger.Info("SERVICE: UPDATE USER")
	_, err := s.storage.FindRecordById(record.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Errorf("failed to get record: %v", err)
		}
		return err
	}
	err = s.storage.UpdateRecord(record)
	if err != nil {
		s.logger.Errorf("failed to update record: %v", err)
		return err
	}
	return nil
}
func (s *service) PartiallyUpdate(ctx context.Context, record *PartiallyUpdateRecordDTO) error {
	s.logger.Info("SERVICE: PARTIALLY UPDATE RECORD")
	_, err := s.storage.FindRecordById(record.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Errorf("failed to get record: %v", err)
		}
		return err
	}
	err = s.storage.PartiallyUpdateRecord(record)
	if err != nil {
		s.logger.Errorf("failed to partially update record: %v", err)
		return err
	}
	return nil
}

func (s *service) Delete(id int64) error {
	s.logger.Info("SERVICE: DELETE RECORD")
	err := s.storage.DeleteRecord(id)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Warnf("failed to delete record: %v", err)
		}
		return err
	}
	return nil
}
