package disease

import (
	"HospitalRecord/app/internal/domain/apperror"
	"HospitalRecord/app/pkg/logger"
	"context"
	"errors"
)

/// Интерфейс Service реализизирующий service и методы для обработки CRUD системы болезей \\\

type Service interface {
	Create(ctx context.Context, input *CreateDiseaseDTO) (*Disease, error)
	GetById(ctx context.Context, id int64) (*Disease, error)
	Update(ctx context.Context, disease *UpdateDiseaseDTO) error
	Delete(id int64) error
}

/// Структура  service реализизирующая инфтерфейс Service болезней \\\

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

/// Функция Create создает болезнь через интерфейс Service принимая входные данные input \\\

func (s *service) Create(ctx context.Context, input *CreateDiseaseDTO) (*Disease, error) {
	s.logger.Info("SERVICE: CREATE DISEASE")

	/// Создание структуры dis на основе полученных данных \\\
	dis := Disease{
		BodyPart:    input.BodyPart,
		Description: input.Description,
	}

	/// Вызов функции Create в хранилище болезней \\\
	disease, err := s.storage.Create(&dis)
	if err != nil {
		return nil, err
	}

	return disease, nil
}

/// Функция GetById осуществялет поиск болезни через интерфейс Service принимая входные данные id \\\

func (s *service) GetById(ctx context.Context, id int64) (*Disease, error) {
	s.logger.Info("SERVICE: GET DISEASE BY ID")

	/// Вызов функции FindById в хранилище болезней \\\
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

/// Функция Update обновляет болезнь через интерфейс Service принимая входные данные disease \\\

func (s *service) Update(ctx context.Context, disease *UpdateDiseaseDTO) error {
	s.logger.Info("SERVICE: UPDATE DISEASE")

	/// Проверка на существование болезни с данным id \\\
	/// Вызов функции FindById в хранилище болезней \\\
	_, err := s.storage.FindById(disease.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Errorf("failed to get disease: %v", err)
		}
		return err
	}

	/// Вызов функции Update в хранилище болезней \\\
	err = s.storage.Update(disease)
	if err != nil {
		s.logger.Errorf("failed to update disease: %v", err)
		return err
	}
	return nil
}

/// Функция Delete удаляет болезнь через интерфейс Service принимая входные данные id \\\

func (s *service) Delete(id int64) error {
	s.logger.Info("SERVICE: DELETE DISEASE")

	/// Вызов функции Delete в хранилище болезней \\\
	err := s.storage.Delete(id)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Warnf("failed to delete disease: %v", err)
		}
		return err
	}
	return nil
}
