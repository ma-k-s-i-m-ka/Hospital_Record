package user

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
	usersURL           = "/hospital_record/users"
	userByEmailURL     = "/hospital_record/users/email"
	userByPolicyNumber = "/hospital_record/users/policy_number"
	userURL            = "/hospital_record/users/profile/:id"
)

/// Структура Handler представляющая собой обработчик объекта userService для пациентов \\\

type Handler struct {
	logger      logger.Logger
	userService Service
}

/// Структура NewHandler возвращает новый экземпляр Handler инициализируя переданные в него аргументы \\\

func NewHandler(logger logger.Logger, userService Service) handler.Hand {
	return &Handler{
		logger:      logger,
		userService: userService,
	}
}

/// Структура Register регистрирует новые запросы для пациентов \\\

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, userByEmailURL, h.GetUserByEmail)
	router.HandlerFunc(http.MethodGet, userByPolicyNumber, h.GetUserByPolicyNumber)
	router.HandlerFunc(http.MethodPost, usersURL, h.CreateUser)
	router.HandlerFunc(http.MethodPut, userURL, h.UpdateUser)
	router.HandlerFunc(http.MethodPatch, userURL, h.PartiallyUpdateUser)
	router.HandlerFunc(http.MethodDelete, userURL, h.DeleteUser)
	router.HandlerFunc(http.MethodGet, userURL, h.GetUserById)
}

/// Функция GetUserById получает пациента по его id \\\

func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: GET USER BY ID")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	h.logger.Printf("Input: %+v\n", id)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	/// Вызов функции GetById передавая ей id пациента \\\
	user, err := h.userService.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		h.logger.Error(err)
		response.InternalError(w, err.Error(), "")
		return
	}
	h.logger.Info("GOT USER BY ID")
	response.JSON(w, http.StatusOK, user)
}

/// Функция GetUserByEmail получает пациента по его email \\\

func (h *Handler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: GET USER BY EMAIL")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр email из URL \\\
	email := r.URL.Query().Get("email")
	h.logger.Printf("Input: %+v\n", email)
	if email == "" {
		response.BadRequest(w, "empty email", "")
		return
	}

	/// Вызов функции GetByEmail передавая ей email \\\
	user, err := h.userService.GetByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		response.BadRequest(w, err.Error(), "")
		return
	}
	h.logger.Info("GOT USER BY EMAIL")
	response.JSON(w, http.StatusOK, user)
}

/// Функция GetUserByPolicyNumber получает пациента по его номеру полиса \\\

func (h *Handler) GetUserByPolicyNumber(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: GET USER BY POLICY NUMBER")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр policy из URL \\\
	policy := r.URL.Query().Get("policy_number")
	if policy == "" {
		response.BadRequest(w, "empty policy", "")
		return
	}

	/// Вызов функции GetByPolicyNumber передавая ей номер полиса \\\
	user, err := h.userService.GetByPolicyNumber(r.Context(), policy)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		response.BadRequest(w, err.Error(), "")
		return
	}
	h.logger.Info("GOT USER BY POLICY NUMBER")
	response.JSON(w, http.StatusOK, user)
}

/// Функция CreateUser создает пациента по полученным данным из input \\\

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: CREATE USER")

	var input CreateUserDTO

	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.logger.Printf("Input: %+v\n", &input)

	/// Вызов функции Create передавая ей полученные значения и ссылку на структуру input \\\
	user, err := h.userService.Create(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrRepeatedEmail) {
			response.BadRequest(w, err.Error(), "")
			return
		}
		response.InternalError(w, fmt.Sprintf("cannot create user: %v", err), "")
		return
	}
	h.logger.Info("USER CREATED")
	response.JSON(w, http.StatusCreated, user)
}

/// Функция UpdateUser обновляет пациента по его id и полученным данным из input \\\

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: UPDATE USER")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
	var input UpdateUserDTO
	input.ID = id

	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.logger.Printf("Input: %+v\n", &input)

	/// Вызов функции Update передавая ей полученные значения и ссылку на структуру input \\\
	err = h.userService.Update(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
	}
	h.logger.Info("USER UPDATED")
	response.JSON(w, http.StatusOK, "USER UPDATED")
}

/// Функция PartiallyUpdateUser частично обновляет пациента по его id и полученным данным из input \\\

func (h *Handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: PARTIALLY UPDATE USER")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
	var input PartiallyUpdateUserDTO
	input.ID = id

	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}

	/// Вызов функции PartiallyUpdate передавая ей полученные значения и ссылку на структуру input \\\
	err = h.userService.PartiallyUpdate(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
	}
	h.logger.Info("USER PARTIALLY UPDATED")
	response.JSON(w, http.StatusOK, "USER PARTIALLY UPDATED")
}

/// Функция DeleteUser удаляет пациента по его id \\\

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: DELETE USER")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	/// Вызов функции Delete передавая ей полученное значение id \\\
	err = h.userService.Delete(id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		response.InternalError(w, err.Error(), "wrong on the server")
		return
	}
	h.logger.Info("USER DELETED")
	response.JSON(w, http.StatusOK, "USER DELETED")
}
