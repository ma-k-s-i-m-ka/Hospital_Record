package record

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
	recordsURL         = "/hospital_record/records"
	recordByPatientsId = "/hospital_record/record/patients_record/:id"
	recordURL          = "/hospital_record/records/:id"
)

/// Структура Handler представляющая собой обработчик объекта recordService для записей на прием \\\

type Handler struct {
	logger        logger.Logger
	recordService Service
}

/// Структура NewHandler возвращает новый экземпляр Handler инициализируя переданные в него аргументы \\\

func NewHandler(logger logger.Logger, recordService Service) handler.Hand {
	return &Handler{
		logger:        logger,
		recordService: recordService,
	}
}

/// Структура Register регистрирует новые запросы для записей на прием \\\

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, recordsURL, h.CreateRecord)
	router.HandlerFunc(http.MethodGet, recordByPatientsId, h.GetRecordByPatientsId)
	router.HandlerFunc(http.MethodPut, recordURL, h.UpdateRecord)
	router.HandlerFunc(http.MethodPatch, recordURL, h.PartiallyUpdateRecord)
	router.HandlerFunc(http.MethodDelete, recordURL, h.DeleteRecord)
	router.HandlerFunc(http.MethodGet, recordURL, h.GetRecordById)
}

/// Функция GetRecordById получает запись по ее id \\\

func (h *Handler) GetRecordById(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: GET RECORD BY ID")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	h.logger.Printf("Input: %+v\n", id)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	/// Вызов функции GetById передавая ей id записи \\\
	record, err := h.recordService.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		h.logger.Error(err)
		response.InternalError(w, err.Error(), "")
		return
	}
	h.logger.Info("GOT RECORD BY ID")
	response.JSON(w, http.StatusOK, record)
}

/// Функция GetRecordByPatientsId получает запись на прием по id пациента \\\

func (h *Handler) GetRecordByPatientsId(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: GET RECORD BY PATIENTS ID")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	h.logger.Printf("Input: %+v\n", id)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	/// Вызов функции GetByPatientsId передавая ей id пациента \\\
	record, err := h.recordService.GetByPatientsId(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		response.BadRequest(w, err.Error(), "")
		return
	}
	h.logger.Info("GOT RECORD BY PATIENTS ID")
	response.JSON(w, http.StatusOK, record)
}

/// Функция CreateRecord создает запись на прием по полученным данным из input \\\

func (h *Handler) CreateRecord(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: CREATE RECORD")
	var input CreateRecordDTO

	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.logger.Printf("Input: %+v\n", &input)

	/// Вызов функции Create передавая ей полученные значения и ссылку на структуру input \\\
	record, err := h.recordService.Create(r.Context(), &input)
	if err != nil {
		response.InternalError(w, fmt.Sprintf("cannot create record: %v", err), "")
		return
	}
	h.logger.Info("RECORD CREATED")
	response.JSON(w, http.StatusCreated, record)
}

/// Функция UpdateRecord обновляет запись на прием по ее id и полученным данным из input \\\

func (h *Handler) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: UPDATE RECORD")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
	var input UpdateRecordDTO
	input.ID = id

	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.logger.Printf("Input: %+v\n", &input)

	/// Вызов функции Update передавая ей полученные значения и ссылку на структуру input \\\
	err = h.recordService.Update(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
	}
	h.logger.Info("RECORD UPDATED")
	response.JSON(w, http.StatusOK, "RECORD UPDATED")
}

/// Функция PartiallyUpdateRecord частично обновляет запись на прием по ее id и полученным данным из input \\\

func (h *Handler) PartiallyUpdateRecord(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: PARTIALLY UPDATE RECORD")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
	var input PartiallyUpdateRecordDTO
	input.ID = id

	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}

	/// Вызов функции PartiallyUpdate передавая ей полученные значения и ссылку на структуру input \\\
	err = h.recordService.PartiallyUpdate(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
	}
	h.logger.Info("RECORD PARTIALLY UPDATED")
	response.JSON(w, http.StatusOK, "RECORD PARTIALLY UPDATED")
}

/// Функция DeleteRecord удаляет запись на прием по ее id \\\

func (h *Handler) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: DELETE RECORD")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	/// Вызов функции Delete передавая ей полученное значение id \\\
	err = h.recordService.Delete(id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		response.InternalError(w, err.Error(), "wrong on the server")
		return
	}
	h.logger.Info("RECORD DELETED")
	response.JSON(w, http.StatusOK, "RECORD DELETED")
}
