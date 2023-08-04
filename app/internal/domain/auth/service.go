package auth

import (
	"HospitalRecord/app/internal/config"
	"HospitalRecord/app/internal/domain/apperror"
	"HospitalRecord/app/internal/domain/user"
	"HospitalRecord/app/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

/// Интерфейс Service реализизирующий service и методы для обработки логики аутентификации и регистрации пользователей \\\

type Service interface {
	AuthByEmail(ctx context.Context, user *AuthByEmail) (*AuthResponse, error)
	AuthByPolicyNumber(ctx context.Context, user *AuthByPolicyNumber) (*AuthResponse, error)
	Register(ctx context.Context, user *Register) (*RegisterResponse, error)
	CreateAccessToken(cfg *config.Config, user *user.User) (string, error)
	CreateRefreshToken(cfg *config.Config, user *user.User) (string, error)
}

/// Структура  service реализизирующая инфтерфейс Service пользователей \\\

type service struct {
	logger  logger.Logger
	storage user.Storage
	cfg     *config.Config
}

/// Структура NewService возвращает новый экземпляр Service инициализируя переданные в него аргументы \\\

func NewService(storage user.Storage, logger logger.Logger, cfg *config.Config) Service {
	return &service{
		logger:  logger,
		storage: storage,
		cfg:     cfg,
	}
}

/// Функция AuthByEmail реализует аутентификацию пользователя по адресу электронной почты через интерфейс Service принимая входные данные input  \\\

func (s *service) AuthByEmail(ctx context.Context, input *AuthByEmail) (*AuthResponse, error) {
	s.logger.Info("SERVICE: AUTH USER BY EMAIL")

	/// Вызов функции FindByEmail в хранилище пользователей  \\\
	user, err := s.storage.FindByEmail(input.Email)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.logger.Warnf("cannot find user by email: %v", err)
		return nil, err
	}
	/// Проверка на соответствие введенного и захэшированного пароля в хранилище \\\
	if !user.CheckPassword(input.Password) {
		s.logger.Warnf("incorrect password: %v", err)
		return nil, err
	}

	/// Создание токенов доступа \\\
	accessToken, err := s.CreateAccessToken(s.cfg, user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.CreateRefreshToken(s.cfg, user)
	if err != nil {
		return nil, err
	}
	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

/// Функция AuthByPolicyNumber реализует аутентификацию пользователя по номеру полиса через интерфейс Service принимая входные данные input  \\\

func (s *service) AuthByPolicyNumber(ctx context.Context, input *AuthByPolicyNumber) (*AuthResponse, error) {
	s.logger.Info("SERVICE: AUTH USER BY POLICY NUMBER")

	/// Вызов функции FindByPolicyNumber в хранилище пользователей  \\\
	user, err := s.storage.FindByPolicyNumber(input.PolicyNumber)

	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.logger.Warnf("cannot find user by policy number: %v", err)
		return nil, err
	}

	/// Проверка на соответствие введенного и захэшированного пароля в хранилище \\\
	if !user.CheckPassword(input.Password) {
		s.logger.Warnf("incorrect password")
		return nil, apperror.ErrInvalidCredentials
	}

	/// Создание токенов доступа \\\
	accessToken, err := s.CreateAccessToken(s.cfg, user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.CreateRefreshToken(s.cfg, user)
	if err != nil {
		return nil, err
	}
	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

/// Функция Register реализует регистрацию пользователя через интерфейс Service принимая входные данные input  \\\

func (s *service) Register(ctx context.Context, input *Register) (*RegisterResponse, error) {
	s.logger.Info("SERVICE: REGISTER USER")

	/// Проверка на повтаряющийся адрес электронной почты \\\
	/// Вызов функции FindByEmail в хранилище пользователей  \\\
	checkEmail, err := s.storage.FindByEmail(input.Email)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return nil, err
		}
	}

	if checkEmail != nil {
		return nil, apperror.ErrRepeatedEmail
	}

	/// Проверка на повтаряющийся номер полиса \\\
	/// Вызов функции FindByPolicyNumber в хранилище пользователей  \\\
	checkPolicyNumber, err := s.storage.FindByPolicyNumber(input.PolicyNumber)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return nil, err
		}
	}

	if checkPolicyNumber != nil {
		return nil, apperror.ErrRepeatedPolicyNumber
	}

	u := user.User{
		Email:        input.Email,
		Name:         input.Name,
		Surname:      input.Surname,
		Age:          input.Age,
		Gender:       input.Gender,
		Password:     input.Password,
		PolicyNumber: input.PolicyNumber,
	}

	/// Хэширование полученного пароля \\\
	err = u.HashPassword()
	if err != nil {
		return nil, fmt.Errorf("cannot hash password")
	}

	/// Вызов функции Create в хранилище пользователей  \\\
	user, err := s.storage.Create(&u)

	/// Создание токенов доступа \\\
	accessToken, err := s.CreateAccessToken(s.cfg, user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.CreateRefreshToken(s.cfg, user)
	if err != nil {
		return nil, err
	}
	return &RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

/// Функция CreateAccessToken для создания токена доступа AccessToken \\\

func (s *service) CreateAccessToken(cfg *config.Config, user *user.User) (string, error) {
	s.logger.Info("SERVICE: CREATE ACCESS TOKEN")
	metadata := AccessToken{
		ID:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		Surname:      user.Surname,
		Patronymic:   user.Patronymic,
		Age:          user.Age,
		Gender:       user.Gender,
		PhoneNumber:  user.PhoneNumber,
		Address:      user.Address,
		PolicyNumber: user.PolicyNumber,
	}

	/// Добавление информации о пользователе в мапу claims, и время истечения действия токена AccessToken \\\
	claims := jwt.MapClaims{
		"user": metadata,
		"exp":  time.Duration(cfg.JWT.AccessExpirationMinutes) * time.Minute,
	}
	/// Создание нового токена accessToken с указанными утверждениями claims и методом подписи SigningMethodHS256 \\\
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	/// Токен подписывается с помощью секретного ключа AccessTokenSecretKey и преобразуется в строку с помощью метода SignedString \\\
	token, err := accessToken.SignedString([]byte(cfg.JWT.AccessTokenSecretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

/// Функция CreateRefreshToken для создания токена обновления RefreshToken \\\

func (s *service) CreateRefreshToken(cfg *config.Config, user *user.User) (string, error) {
	s.logger.Info("SERVICE: CREATE REFRESH TOKEN")

	/// Добавление id пользователя в мапу claims, и время истечения действия токена RefreshToken \\\
	claims := jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Duration(cfg.JWT.RefreshExpirationDays) * time.Hour * 24,
	}
	/// Создание нового токена refreshToken с указанными утверждениями claims и методом подписи SigningMethodHS256 \\\
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	/// Токен подписывается с помощью секретного ключа RefreshTokenSecretKey и преобразуется в строку с помощью метода SignedString \\\
	token, err := refreshToken.SignedString([]byte(cfg.JWT.RefreshTokenSecretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}
