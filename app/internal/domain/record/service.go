package record

import (
	"HospitalRecord/app/internal/domain/apperror"
	"HospitalRecord/app/internal/domain/doctor"
	"HospitalRecord/app/pkg/logger"
	"context"
	"errors"
)

/// Интерфейс Service реализизирующий service и методы для обработки CRUD системы записей на прием \\\

type Service interface {
	Create(ctx context.Context, record *CreateRecordDTO) (*Record, error)
	GetByPatientsId(ctx context.Context, id int64) (*Record, error)
	GetById(ctx context.Context, id int64) (*Record, error)
	Update(ctx context.Context, record *UpdateRecordDTO) error
	PartiallyUpdate(ctx context.Context, record *PartiallyUpdateRecordDTO) error
	Delete(id int64) error
}

/// Структура  service реализизирующая инфтерфейс Service записей на прием \\\

type service struct {
	logger  logger.Logger
	storage Storage
	doc     doctor.Storage
}

/// Структура NewService возвращает новый экземпляр Service инициализируя переданные в него аргументы \\\

func NewService(doc doctor.Storage, storage Storage, logger logger.Logger) Service {
	return &service{
		logger:  logger,
		storage: storage,
		doc:     doc,
	}
}

/// Функция Create создает запись через интерфейс Service принимая входные данные input \\\

func (s *service) Create(ctx context.Context, input *CreateRecordDTO) (*Record, error) {
	s.logger.Info("SERVICE: CREATE RECORD")

	/// Проверка что доктор свободен \\\
	checkDoctor, err := s.doc.FindById(input.DoctorID)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return nil, err
		}
	}
	if checkDoctor.RecordingIsAvailable != true {
		return nil, apperror.ErrDoctorNotAvailable
	}

	/// Создание структуры r на основе полученных данных \\\
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

	/// Вызов функции Create в хранилище записей \\\
	record, err := s.storage.CreateRecord(&r)
	if err != nil {
		return nil, err
	}
	return record, nil
}

/// Функция GetByPatientsId осуществялет поиск записи на прием через интерфейс Service принимая входные данные id пациента\\\

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

/// Функция GetById осуществялет поиск записи на прием через интерфейс Service принимая входные данные id \\\

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

/// Функция Update обновляет запись через интерфейс Service принимая входные данные record \\\

func (s *service) Update(ctx context.Context, record *UpdateRecordDTO) error {
	s.logger.Info("SERVICE: UPDATE USER")

	/// Вызов функции FindRecordById в хранилище записей \\\
	_, err := s.storage.FindRecordById(record.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Errorf("failed to get record: %v", err)
		}
		return err
	}

	/// Вызов функции UpdateRecord в хранилище записей \\\
	err = s.storage.UpdateRecord(record)
	if err != nil {
		s.logger.Errorf("failed to update record: %v", err)
		return err
	}
	return nil
}

/// Функция PartiallyUpdate частично обновляет запись на прием через интерфейс Service принимая входные данные record \\\

func (s *service) PartiallyUpdate(ctx context.Context, record *PartiallyUpdateRecordDTO) error {
	s.logger.Info("SERVICE: PARTIALLY UPDATE RECORD")

	/// Вызов функции FindRecordById в хранилище записей \\\
	_, err := s.storage.FindRecordById(record.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Errorf("failed to get record: %v", err)
		}
		return err
	}

	/// Вызов функции PartiallyUpdateRecord в хранилище записей \\\
	err = s.storage.PartiallyUpdateRecord(record)
	if err != nil {
		s.logger.Errorf("failed to partially update record: %v", err)
		return err
	}
	return nil
}

/// Функция Delete удаляет запись на прием через интерфейс Service принимая входные данные id \\\

func (s *service) Delete(id int64) error {
	s.logger.Info("SERVICE: DELETE RECORD")

	/// Вызов функции Delete в хранилище записей \\\
	err := s.storage.DeleteRecord(id)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Warnf("failed to delete record: %v", err)
		}
		return err
	}
	return nil
}
