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

// Handler handles requests specified to user service.
type Handler struct {
	logger        logger.Logger
	recordService Service
}

// NewHandler returns a new user Handler instance.
func NewHandler(logger logger.Logger, recordService Service) handler.Hand {
	return &Handler{
		logger:        logger,
		recordService: recordService,
	}
}

// Register registers new routes for router.
func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, recordsURL, h.CreateRecord)
	router.HandlerFunc(http.MethodGet, recordByPatientsId, h.GetRecordByPatientsId)
	router.HandlerFunc(http.MethodPut, recordURL, h.UpdateRecord)
	router.HandlerFunc(http.MethodPatch, recordURL, h.PartiallyUpdateRecord)
	router.HandlerFunc(http.MethodDelete, recordURL, h.DeleteRecord)
	router.HandlerFunc(http.MethodGet, recordURL, h.GetRecordById)
}

func (h *Handler) GetRecordById(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: GET RECORD BY ID")

	id, err := handler.ReadIdParam64(r)
	h.logger.Printf("Input: %+v\n", id)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
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
	response.JSON(w, http.StatusOK, record)
}

func (h *Handler) GetRecordByPatientsId(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: GET RECORD BY PATIENTS ID")

	id, err := handler.ReadIdParam64(r)
	h.logger.Printf("Input: %+v\n", id)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
	record, err := h.recordService.GetByPatientsId(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		response.BadRequest(w, err.Error(), "")
		return
	}
	response.JSON(w, http.StatusOK, record)
}

func (h *Handler) CreateRecord(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: CREATE RECORD")

	var input CreateRecordDTO
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.logger.Printf("Input: %+v\n", &input)
	record, err := h.recordService.Create(r.Context(), &input)
	if err != nil {
		response.InternalError(w, fmt.Sprintf("cannot create record: %v", err), "")
		return
	}

	response.JSON(w, http.StatusCreated, record)
}

func (h *Handler) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: UPDATE RECORD")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
	var input UpdateRecordDTO
	input.ID = id
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.logger.Printf("Input: %+v\n", &input)
	err = h.recordService.Update(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) PartiallyUpdateRecord(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: PARTIALLY UPDATE RECORD")
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
	var input PartiallyUpdateRecordDTO
	input.ID = id
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	err = h.recordService.PartiallyUpdate(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: DELETE RECORD BY ID")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	err = h.recordService.Delete(id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		response.InternalError(w, err.Error(), "wrong on the server")
		return
	}

	w.WriteHeader(http.StatusOK)
}
