package user

import (
	"HospitalRecord/app/internal/domain/apperror"
	"HospitalRecord/app/pkg/logger"
	"context"
	"errors"
	"fmt"
)

/// Интерфейс Service реализизирующий service и методы для обработки CRUD системы пациентов \\\

type Service interface {
	Create(ctx context.Context, user *CreateUserDTO) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetById(ctx context.Context, id int64) (*User, error)
	GetByPolicyNumber(ctx context.Context, policy string) (*User, error)
	Update(ctx context.Context, user *UpdateUserDTO) error
	PartiallyUpdate(ctx context.Context, user *PartiallyUpdateUserDTO) error
	Delete(id int64) error
}

/// Структура  service реализизирующая инфтерфейс Service пациентов \\\

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

/// Функция Create создает пациента через интерфейс Service принимая входные данные input \\\

func (s *service) Create(ctx context.Context, input *CreateUserDTO) (*User, error) {
	s.logger.Info("SERVICE: CREATE USER")

	/// Проверка на уникальность email \\\
	checkEmail, err := s.storage.FindByEmail(input.Email)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return nil, err
		}
	}
	if checkEmail != nil {
		return nil, apperror.ErrRepeatedEmail
	}

	/// Проверка на уникальность номера полиса \\\
	checkPolicyNumber, err := s.storage.FindByPolicyNumber(input.PolicyNumber)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return nil, err
		}
	}
	if checkPolicyNumber != nil {
		return nil, apperror.ErrRepeatedPolicyNumber
	}

	/// Создание структуры u на основе полученных данных \\\
	u := User{
		Email:        input.Email,
		Name:         input.Name,
		Surname:      input.Surname,
		Age:          input.Age,
		Gender:       input.Gender,
		Password:     input.Password,
		PolicyNumber: input.PolicyNumber,
	}

	/// Хэширование пароля \\\
	err = u.HashPassword()
	if err != nil {
		return nil, fmt.Errorf("cannot hash password")
	}

	/// Вызов функции Create в хранилище пациентов \\\
	user, err := s.storage.Create(&u)
	if err != nil {
		return nil, err
	}
	return user, nil
}

/// Функция GetByEmail осуществялет поиск пациентов через интерфейс Service принимая входные данные email пациента \\\

func (s *service) GetByEmail(ctx context.Context, email string) (*User, error) {
	s.logger.Info("SERVICE: GET USER BY EMAIL")

	/// Вызов функции FindByEmail в хранилище пациентов \\\
	user, err := s.storage.FindByEmail(email)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.logger.Warnf("cannot find user by email: %v", err)
		return nil, err
	}
	return user, nil
}

/// Функция GetById осуществялет поиск пациентов через интерфейс Service принимая входные данные id пациента \\\

func (s *service) GetById(ctx context.Context, id int64) (*User, error) {
	s.logger.Info("SERVICE: GET USER BY ID")

	/// Вызов функции FindById в хранилище пациентов \\\
	user, err := s.storage.FindById(id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.logger.Warnf("cannot find user by id: %v", err)
		return nil, err
	}
	return user, nil
}

/// Функция GetByPolicyNumber осуществялет поиск пациентов через интерфейс Service принимая входные данные номер полиса пациента \\\

func (s *service) GetByPolicyNumber(ctx context.Context, policy string) (*User, error) {
	s.logger.Info("SERVICE: GET USER BY POLICY NUMBER")

	/// Вызов функции FindByPolicyNumber в хранилище пациентов \\\
	user, err := s.storage.FindByPolicyNumber(policy)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.logger.Warnf("cannot find user by policy number: %v", err)
		return nil, err
	}
	return user, nil
}

/// Функция Update обновляет пациентов через интерфейс Service принимая входные данные user \\\

func (s *service) Update(ctx context.Context, user *UpdateUserDTO) error {
	s.logger.Info("SERVICE: UPDATE USER")

	/// Вызов функции FindById в хранилище пациентов \\\
	u, err := s.storage.FindById(user.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Errorf("failed to get user: %v", err)
		}
		return err
	}

	/// Хэширование обновленного пароля \\\
	u.Password = user.Password
	err = user.HashPassword()
	if err != nil {
		s.logger.Errorf("failed to hash password: %v", err)
		return err
	}

	/// Вызов функции Update в хранилище пациентов \\\
	err = s.storage.Update(user)
	if err != nil {
		s.logger.Errorf("failed to update user: %v", err)
		return err
	}
	return nil
}

/// Функция PartiallyUpdate частично обновляет пациента через интерфейс Service принимая входные данные user \\\

func (s *service) PartiallyUpdate(ctx context.Context, user *PartiallyUpdateUserDTO) error {
	s.logger.Info("SERVICE: PARTIALLY UPDATE USER")

	/// Вызов функции FindById в хранилище пациентов \\\
	u, err := s.storage.FindById(user.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Errorf("failed to get user: %v", err)
		}
		return err
	}

	/// Хэширование обновленного пароля \\\
	if user.Password != nil {
		u.Password = *user.Password
		err = user.HashPassword()
		if err != nil {
			s.logger.Errorf("failed to hash password: %v", err)
			return err
		}
	}

	/// Вызов функции PartiallyUpdate в хранилище пациентов \\\
	err = s.storage.PartiallyUpdate(user)
	if err != nil {
		s.logger.Errorf("failed to partially update user: %v", err)
		return err
	}
	return nil
}

/// Функция Delete удаляет пациента через интерфейс Service принимая входные данные id \\\

func (s *service) Delete(id int64) error {
	s.logger.Info("SERVICE: DELETE USER")

	/// Вызов функции Delete в хранилище пациентов \\\
	err := s.storage.Delete(id)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.logger.Warnf("failed to delete user: %v", err)
		}
		return err
	}
	return nil
}
