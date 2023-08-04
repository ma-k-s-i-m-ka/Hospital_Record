package specialization

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
	specializationsURL = "/hospital_record/specializations"
	specializationURL  = "/hospital_record/specializations/:id"
)

/// Структура Handler представляющая собой обработчик объекта specializationService для специализаций \\\

type Handler struct {
	logger                logger.Logger
	specializationService Service
}

/// Структура NewHandler возвращает новый экземпляр Handler инициализируя переданные в него аргументы \\\

func NewHandler(logger logger.Logger, specializationService Service) handler.Hand {
	return &Handler{
		logger:                logger,
		specializationService: specializationService,
	}
}

/// Структура Register регистрирует новые запросы для специализации \\\

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, specializationURL, h.GetSpecializationById)
	router.HandlerFunc(http.MethodPost, specializationsURL, h.CreateSpecialization)
	router.HandlerFunc(http.MethodPut, specializationURL, h.UpdateSpecialization)
	router.HandlerFunc(http.MethodDelete, specializationURL, h.DeleteSpecialization)
}

/// Функция GetSpecializationById получает специализацию по ee id \\\

func (h *Handler) GetSpecializationById(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: GET SPECIALIZATION BY ID")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	/// Вызов функции GetById передавая ей id специализации \\\
	specialization, err := h.specializationService.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		h.logger.Error(err)
		response.InternalError(w, err.Error(), "")
		return
	}
	h.logger.Info("GOT SPECIALIZATION BY ID")
	response.JSON(w, http.StatusOK, specialization)
}

/// Функция CreateSpecialization создает специализацию по полученным данным из input \\\

func (h *Handler) CreateSpecialization(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: CREATE SPECIALIZATION")

	var input CreateSpecializationDTO

	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.logger.Printf("Input: %+v\n", &input)

	/// Вызов функции Create передавая ей полученные значения и ссылку на структуру input \\\
	specialization, err := h.specializationService.Create(r.Context(), &input)
	if err != nil {
		response.InternalError(w, fmt.Sprintf("cannot create specialization: %v", err), "")
		return
	}
	h.logger.Info("SPECIALIZATION CREATED")
	response.JSON(w, http.StatusCreated, specialization)
}

/// Функция UpdateSpecialization обновляет специализацию по ее id и полученным данным из input \\\

func (h *Handler) UpdateSpecialization(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: UPDATE SPECIALIZATION")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	var input UpdateSpecializationDTO
	input.ID = id

	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.logger.Printf("Input: %+v\n", &input)

	/// Вызов функции Update передавая ей полученные значения и ссылку на структуру input \\\
	err = h.specializationService.Update(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
	}
	h.logger.Info("SPECIALIZATION UPDATED")
	response.JSON(w, http.StatusOK, "SPECIALIZATION UPDATED")
}

/// Функция DeleteSpecialization удаляет специализацию по ее id \\\

func (h *Handler) DeleteSpecialization(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: DELETE SPECIALIZATION")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	/// Вызов функции Delete передавая ей полученное значение id \\\
	err = h.specializationService.Delete(id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		response.InternalError(w, err.Error(), "wrong on the server")
		return
	}
	h.logger.Info("SPECIALIZATION DELETED")
	response.JSON(w, http.StatusOK, "SPECIALIZATION DELETED")
}
