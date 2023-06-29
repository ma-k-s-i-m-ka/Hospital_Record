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
	specializationsURL = "/hospitalrecord/specializations"
	specializationURL  = "/hospitalrecord/specializations/:id"
)

// Handler handles requests specified to user service.
type Handler struct {
	logger                logger.Logger
	specializationService Service
}

// NewHandler returns a new user Handler instance.
func NewHandler(logger logger.Logger, specializationService Service) handler.Hand {
	return &Handler{
		logger:                logger,
		specializationService: specializationService,
	}
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, specializationURL, h.GetSpecializationById)
	router.HandlerFunc(http.MethodPost, specializationsURL, h.CreateSpecialization)
	router.HandlerFunc(http.MethodPut, specializationURL, h.UpdateSpecialization)
	router.HandlerFunc(http.MethodDelete, specializationURL, h.DeleteSpecialization)
}

func (h *Handler) GetSpecializationById(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("GET SPECIALIZATION BY ID")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
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

	response.JSON(w, http.StatusOK, specialization)
}

func (h *Handler) CreateSpecialization(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("CREATE SPECIALIZATION")

	var input CreateSpecializationDTO

	specialization, err := h.specializationService.Create(r.Context(), &input)
	if err != nil {
		response.InternalError(w, fmt.Sprintf("cannot create specialization: %v", err), "")
		return
	}

	response.JSON(w, http.StatusCreated, specialization)
}

func (h *Handler) UpdateSpecialization(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("UPDATE SPECIALIZATION")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	var input UpdateSpecializationDTO

	input.ID = id

	err = h.specializationService.Update(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteSpecialization(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("DELETE SPECIALIZATION")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	err = h.specializationService.Delete(r.Context(), id)
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
