package doctor

import (
	"HospitalRecord/app/internal/domain/apperror"
	"HospitalRecord/app/pkg/logger"
	"context"
	"errors"
)

/// Интерфейс Service реализизирующий service и методы для обработки CRUD системы докторов \\\

type Service interface {
	Create(ctx context.Context, doctor *CreateDoctorDTO) (*Doctor, error)
	FindAll() (*[]Doctor, error)
	FindAllAvailable(ctx context.Context, id int64, recordingIsAvailable bool) (*[]Doctor, error)
	GetById(ctx context.Context, id int64) (*Doctor, error)
	GetByPortfolioId(ctx context.Context, id int64) (*Doctor, error)
	Update(ctx context.Context, doctor *UpdateDoctorDTO) error
	PartiallyUpdate(ctx context.Context, doctor *PartiallyUpdateDoctorDTO) error
	Delete(id int64) error
}

/// Структура  service реализизирующая инфтерфейс Service докторов \\\

type service struct {
	logger  logger.Logger
	storage Storage
}

/// Структура NewService возвращает новый экземпляр Service инициализируя переданные в него аргументы \\\

func NewService(storage Storage, logger logger.Logger) Service {
	return &service{
		logger:  logger,
		storage: storage,
	}
}

/// Функция Create создает доктора через интерфейс Service принимая входные данные input \\\

func (s *service) Create(ctx context.Context, input *CreateDoctorDTO) (*Doctor, error) {
	s.logger.Info("SERVICE: CREATE DOCTOR")

	/// Проверка на уникальность портфолио \\\
	checkPortfolioID, err := s.storage.FindByPortfolioId(input.PortfolioID)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return nil, err
		}
	}

	if checkPortfolioID != nil {
		return nil, apperror.ErrRepeatedPortfolioId
	}

	/// Создание структуры doc на основе полученных данных \\\
	doc := Doctor{
		Name:                 input.Name,
		Surname:              input.Surname,
		ImageID:              input.ImageID,
		Gender:               input.Gender,
		Rating:               input.Rating,
		Age:                  input.Age,
		RecordingIsAvailable: input.RecordingIsAvailable,
		SpecializationID:     input.SpecializationID,
		PortfolioID:          input.PortfolioID,
	}

	/// Вызов функции Create в хранилище докторов \\\
	doctor, err := s.storage.Create(&doc)
	if err != nil {
		return nil, err
	}
	return doctor, nil
}

/// Функция FindAll осуществялет поиск всех докторов \\\

func (s *service) FindAll() (*[]Doctor, error) {
	s.logger.Info("SERVICE: GET ALL DOCTORS")

	/// Вызов функции FindAll в хранилище докторов \\\
	doctor, err := s.storage.FindAll()
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.logger.Warnf("cannot find doctors by id: %v", err)
		return nil, err
	}
	return &doctor, nil
}

/// Функция FindAllAvailable осуществялет поиск всех свободных докторов по id специализации \\\

func (s *service) FindAllAvailable(ctx context.Context, id int64, recordingIsAvailable bool) (*[]Doctor, error) {
	s.logger.Info("SERVICE: GET ALL AVAILABLE DOCTORS")

	/// Вызов функции FindAllAvailable в хранилище докторов \\\
	doctor, err := s.storage.FindAllAvailable(id, recordingIsAvailable)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.logger.Warnf("cannot find available doctor by id: %v", err)
		return nil, err
	}
	return &doctor, nil
}

/// Функция GetById осуществялет поиск доктора через интерфейс Service принимая входные данные id доктора \\\

func (s *service) GetById(ctx context.Context, id int64) (*Doctor, error) {
	s.logger.Info("SERVICE: GET DOCTOR BY ID")
	s.logger.Printf("Input: %+v\n", id)

	/// Вызов функции FindById в хранилище докторов \\\
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

/// Функция GetByPortfolioId осуществялет поиск доктора через интерфейс Service принимая входные данные id портфолто \\\

func (s *service) GetByPortfolioId(ctx context.Context, id int64) (*Doctor, error) {
	s.logger.Info("SERVICE: GET DOCTOR BY PORTFOLIO ID")

	/// Вызов функции FindByPortfolioId в хранилище докторов \\\
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

/// Функция Update обновляет доктора через интерфейс Service принимая входные данные doctor \\\

func (s *service) Update(ctx context.Context, doctor *UpdateDoctorDTO) error {
	s.logger.Info("SERVICE: UPDATE DOCTOR")

	/// Вызов функции FindById в хранилище докторов \\\
	_, err := s.storage.FindById(doctor.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Errorf("failed to get doctor: %v", err)
		}
		return err
	}

	/// Вызов функции Update в хранилище докторов \\\
	err = s.storage.Update(doctor)
	if err != nil {
		s.logger.Errorf("failed to update doctor: %v", err)
		return err
	}
	return nil
}

/// Функция PartiallyUpdate частично обновляет доктора через интерфейс Service принимая входные данные doctor \\\

func (s *service) PartiallyUpdate(ctx context.Context, doctor *PartiallyUpdateDoctorDTO) error {
	s.logger.Info("SERVICE: PARTIALLY UPDATE DOCTOR")

	/// Вызов функции FindById в хранилище докторов \\\
	_, err := s.storage.FindById(doctor.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Errorf("failed to get doctor: %v", err)
		}
		return err
	}

	/// Вызов функции PartiallyUpdate в хранилище докторов \\\
	err = s.storage.PartiallyUpdate(doctor)
	if err != nil {
		s.logger.Errorf("failed to partially update doctor: %v", err)
		return err
	}
	return nil
}

/// Функция Delete удаляет доктора через интерфейс Service принимая входные данные id \\\

func (s *service) Delete(id int64) error {
	s.logger.Info("SERVICE: DELETE DOCTOR")

	/// Вызов функции Delete в хранилище докторов \\\
	err := s.storage.Delete(id)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Warnf("failed to delete doctor: %v", err)
		}
		return err
	}
	return nil
}
