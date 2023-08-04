package auth

import (
	"HospitalRecord/app/internal/domain/apperror"
	"HospitalRecord/app/internal/domain/handler"
	"HospitalRecord/app/internal/domain/response"
	"HospitalRecord/app/pkg/logger"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	userAuthByPoliciURL = "/hospital_record/user/auth/policy"
	userAuthByEmailURL  = "/hospital_record/user/auth/email"
	userRegisterURL     = "/hospital_record/user/sign_up"
)

/// Структура Handler представляющая собой обработчик объекта authService для пользователей \\\

type Handler struct {
	logger      logger.Logger
	authService Service
}

/// Структура NewHandler возвращает новый экземпляр Handler инициализируя переданные в него аргументы \\\

func NewHandler(logger logger.Logger, authService Service) handler.Hand {
	return &Handler{
		logger:      logger,
		authService: authService,
	}
}

/// Структура Register регистрирует новые запросы для авторизации \\\

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, userAuthByPoliciURL, h.GetUserByPolicyNumber)
	router.HandlerFunc(http.MethodPost, userAuthByEmailURL, h.GetUserByEmail)
	router.HandlerFunc(http.MethodPost, userRegisterURL, h.RegisterUser)
}

/// Функция GetUserByPolicyNumber получает пользователя по его номеру полиса и паролю \\\

func (h *Handler) GetUserByPolicyNumber(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: AUTH BY POLICY NUMBER")

	var input AuthByPolicyNumber
	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.logger.Printf("Input: %+v\n", &input)

	/// Вызов функции AuthByPolicyNumber передавая ей полученные значения и ссылку на структуру input \\\
	user, err := h.authService.AuthByPolicyNumber(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		response.BadRequest(w, err.Error(), "")
		return
	}
	h.logger.Info("AUTH BY POLICY NUMBER IS COMPLETED")
	response.JSON(w, http.StatusOK, user)
}

/// Функция GetUserByEmail получает пользователя по его адресу электронной почты и паролю \\\

func (h *Handler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: AUTH BY EMAIL")

	var input AuthByEmail
	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.logger.Printf("Input: %+v\n", &input)

	/// Вызов функции AuthByEmail передавая ей полученные значения и ссылку на структуру input \\\
	user, err := h.authService.AuthByEmail(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		response.BadRequest(w, err.Error(), "")
		return
	}
	h.logger.Info("AUTH BY EMAIL IS COMPLETED")
	response.JSON(w, http.StatusOK, user)
}

/// Функция RegisterUser регистрирует пользователя \\\

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: REGISTER USER")

	var input Register
	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.logger.Printf("Input: %+v\n", &input)

	/// Вызов функции Register передавая ей полученные значения и ссылку на структуру input \\\
	user, err := h.authService.Register(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrRepeatedEmail) {
			response.BadRequest(w, err.Error(), "")
			return
		}
		response.InternalError(w, fmt.Sprintf("cannot create user: %v", err), "")
		return
	}
	h.logger.Info("REGISTER USER IS COMPLETED")
	response.JSON(w, http.StatusCreated, user)
}
