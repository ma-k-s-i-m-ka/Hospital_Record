package portfolio

import (
	"HospitalRecord/app/internal/domain/apperror"
	"HospitalRecord/app/pkg/logger"
	"context"
	"errors"
)

/// Интерфейс Service реализизирующий service и методы для обработки CRUD системы портфолио докторов \\\

type Service interface {
	Create(ctx context.Context, input *CreatePortfolioDTO) (*Portfolio, error)
	GetById(ctx context.Context, id int64) (*Portfolio, error)
	Update(ctx context.Context, portfolio *UpdatePortfolioDTO) error
	Delete(id int64) error
}

/// Структура  service реализизирующая инфтерфейс Service портфолио \\\

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

/// Функция Create создает портфолио через интерфейс Service принимая входные данные input \\\

func (s *service) Create(ctx context.Context, input *CreatePortfolioDTO) (*Portfolio, error) {
	s.logger.Info("SERVICE: CREATE PORTFOLIO")

	/// Создание структуры portf на основе полученных данных \\\
	portf := Portfolio{
		Education:      input.Education,
		Awards:         input.Awards,
		WorkExperience: input.WorkExperience,
	}

	/// Вызов функции Create в хранилище докторов \\\

	portfolio, err := s.storage.Create(&portf)
	if err != nil {
		return nil, err
	}

	return portfolio, nil
}

/// Функция GetById осуществялет поиск портфолио через интерфейс Service принимая входные данные id портфолио \\\

func (s *service) GetById(ctx context.Context, id int64) (*Portfolio, error) {
	s.logger.Info("SERVICE: GET PORTFOLIO BY ID")
	s.logger.Printf("Input: %+v\n", id)

	/// Вызов функции FindById в хранилище портфолио \\\
	portfolio, err := s.storage.FindById(id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.logger.Warnf("cannot find portfolio by id: %v", err)
		return nil, err
	}
	return portfolio, nil
}

/// Функция Update обновляет портфолио через интерфейс Service принимая входные данные portfolio \\\

func (s *service) Update(ctx context.Context, portfolio *UpdatePortfolioDTO) error {
	s.logger.Info("SERVICE: UPDATE PORTFOLIO")

	/// Вызов функции FindById в хранилище портфолио \\\
	_, err := s.storage.FindById(portfolio.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Errorf("failed to get portfolio: %v", err)
		}
		return err
	}

	/// Вызов функции Update в хранилище портфолио \\\
	err = s.storage.Update(portfolio)
	if err != nil {
		s.logger.Errorf("failed to update portfolio: %v", err)
		return err
	}
	return nil
}

/// Функция Delete удаляет портфолио через интерфейс Service принимая входные данные id \\\

func (s *service) Delete(id int64) error {
	s.logger.Info("SERVICE: DELETE PORTFOLIO")

	/// Вызов функции Delete в хранилище докторов \\\
	err := s.storage.Delete(id)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Warnf("failed to delete portfolio: %v", err)
		}
		return err
	}
	return nil
}
