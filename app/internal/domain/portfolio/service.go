package portfolio

import (
	"HospitalRecord/app/internal/domain/apperror"
	"HospitalRecord/app/pkg/logger"
	"context"
	"errors"
)

type Service interface {
	Create(ctx context.Context, input *CreatePortfolioDTO) (*Portfolio, error)
	GetById(ctx context.Context, id int64) (*Portfolio, error)
	Update(ctx context.Context, portfolio *UpdatePortfolioDTO) error
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

func (s *service) Create(ctx context.Context, input *CreatePortfolioDTO) (*Portfolio, error) {
	s.logger.Info("SERVICE: CREATE PORTFOLIO")
	portf := Portfolio{
		Education:      input.Education,
		Awards:         input.Awards,
		WorkExperience: input.WorkExperience,
	}
	portfolio, err := s.storage.Create(&portf)
	if err != nil {
		return nil, err
	}

	return portfolio, nil
}

func (s *service) GetById(ctx context.Context, id int64) (*Portfolio, error) {
	s.logger.Info("SERVICE: GET PORTFOLIO BY ID")
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

func (s *service) Update(ctx context.Context, portfolio *UpdatePortfolioDTO) error {
	s.logger.Info("SERVICE: UPDATE PORTFOLIO")
	_, err := s.storage.FindById(portfolio.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Errorf("failed to get portfolio: %v", err)
		}
		return err
	}
	err = s.storage.Update(portfolio)
	if err != nil {
		s.logger.Errorf("failed to update portfolio: %v", err)
		return err
	}
	return nil
}

func (s *service) Delete(id int64) error {
	s.logger.Info("SERVICE: DELETE PORTFOLIO")
	err := s.storage.Delete(id)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Warnf("failed to delete portfolio: %v", err)
		}
		return err
	}
	return nil
}
