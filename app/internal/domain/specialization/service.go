package specialization

import (
	"HospitalRecord/app/internal/domain/apperror"
	"HospitalRecord/app/pkg/logger"
	"context"
	"errors"
)

/// Интерфейс Service реализизирующий service и методы для обработки CRUD системы специализации докторов \\\

type Service interface {
	Create(ctx context.Context, input *CreateSpecializationDTO) (*Specialization, error)
	GetById(ctx context.Context, id int64) (*Specialization, error)
	Update(ctx context.Context, specialization *UpdateSpecializationDTO) error
	Delete(id int64) error
}

/// Структура  service реализизирующая инфтерфейс Service специализации докторов \\\

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

/// Функция Create создает специализацию через интерфейс Service принимая входные данные input \\\

func (s *service) Create(ctx context.Context, input *CreateSpecializationDTO) (*Specialization, error) {
	s.logger.Info("SERVICE: CREATE SPECIALIZATION")

	/// Создание структуры specializ на основе полученных данных \\\
	specializ := Specialization{
		Name: input.Name,
	}
	/// Вызов функции Create в хранилище специализации \\\
	specialization, err := s.storage.Create(&specializ)
	if err != nil {
		return nil, err
	}
	return specialization, nil
}

/// Функция GetById осуществялет поиск специализации через интерфейс Service принимая входные данные id \\\

func (s *service) GetById(ctx context.Context, id int64) (*Specialization, error) {
	s.logger.Info("SERVICE: GET SPECIALIZATION BY ID")

	/// Вызов функции FindById в хранилище специализаций \\\
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

/// Функция Update обновляет специализации через интерфейс Service принимая входные данные specialization \\\

func (s *service) Update(ctx context.Context, specialization *UpdateSpecializationDTO) error {
	s.logger.Info("SERVICE: UPDATE SPECIALIZATION")

	/// Вызов функции FindById в хранилище специализаций \\\
	_, err := s.storage.FindById(specialization.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Errorf("failed to get specialization: %v", err)
		}
		return err
	}

	/// Вызов функции Update в хранилище специализаций \\\
	err = s.storage.Update(specialization)
	if err != nil {
		s.logger.Errorf("failed to update specialization: %v", err)
		return err
	}
	return nil
}

/// Функция Delete удаляет специализацию через интерфейс Service принимая входные данные id \\\

func (s *service) Delete(id int64) error {
	s.logger.Info("SERVICE: DELETE SPECIALIZATION")

	/// Вызов функции Delete в хранилище специализаций \\\
	err := s.storage.Delete(id)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Warnf("failed to delete specialization: %v", err)
		}
		return err
	}
	return nil
}
