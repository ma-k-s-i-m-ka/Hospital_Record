package disease

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
	diseasesURL = "/hospital_record/diseases"
	diseaseURL  = "/hospital_record/diseases/:id"
)

/// Структура Handler представляющая собой обработчик объекта diseasesService для болезней \\\

type Handler struct {
	logger          logger.Logger
	diseasesService Service
}

/// Структура NewHandler возвращает новый экземпляр Handler инициализируя переданные в него аргументы \\\

func NewHandler(logger logger.Logger, diseasesService Service) handler.Hand {
	return &Handler{
		logger:          logger,
		diseasesService: diseasesService,
	}
}

/// Структура Register регистрирует новые запросы для болезней \\\

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, diseaseURL, h.GetDiseaseById)
	router.HandlerFunc(http.MethodPost, diseasesURL, h.CreateDisease)
	router.HandlerFunc(http.MethodPut, diseaseURL, h.UpdateDisease)
	router.HandlerFunc(http.MethodDelete, diseaseURL, h.DeleteDisease)
}

/// Функция GetDiseaseById получает болезнь по его id \\\

func (h *Handler) GetDiseaseById(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: GET DISEASE BY ID")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	/// Вызов функции GetById передавая ей полученное значение \\\
	disease, err := h.diseasesService.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		h.logger.Error(err)
		response.InternalError(w, err.Error(), "")
		return
	}
	h.logger.Info("GOT DISEASE BY ID")
	response.JSON(w, http.StatusOK, disease)
}

/// Функция CreateDisease  создает болезнь по полученным данным из input \\\

func (h *Handler) CreateDisease(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: CREATE DISEASE")
	var input CreateDiseaseDTO

	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.logger.Printf("Input: %+v\n", &input)

	/// Вызов функции Create передавая ей полученные значения и ссылку на структуру input \\\
	disease, err := h.diseasesService.Create(r.Context(), &input)
	if err != nil {
		response.InternalError(w, fmt.Sprintf("cannot create disease: %v", err), "")
		return
	}
	h.logger.Info("DISEASE CREATED")
	response.JSON(w, http.StatusCreated, disease)
}

/// Функция UpdateDisease обновляет болезнь по его id и полученным данным из input \\\

func (h *Handler) UpdateDisease(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: UPDATE DISEASE")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
	var input UpdateDiseaseDTO

	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	input.ID = id
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.logger.Printf("Input: %+v\n", &input)

	/// Вызов функции Update передавая ей полученные значения и ссылку на структуру input \\\
	err = h.diseasesService.Update(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
	}
	h.logger.Info("DISEASE UPDATED")
	response.JSON(w, http.StatusOK, "DISEASE UPDATED")
}

/// Функция DeleteDisease удаляет болезнь по его id \\\

func (h *Handler) DeleteDisease(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: DELETE DISEASE")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	/// Вызов функции Delete передавая ей полученное значение id \\\
	err = h.diseasesService.Delete(id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		response.InternalError(w, err.Error(), "wrong on the server")
		return
	}
	h.logger.Info("DISEASE DELETED")
	response.JSON(w, http.StatusOK, "DISEASE DELETED")
}
