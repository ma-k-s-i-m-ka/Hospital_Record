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
	diseasesURL = "/hospitalrecord/diseases"
	diseaseURL  = "/hospitalrecord/diseases/:id"
)

// Handler handles requests specified to user service.
type Handler struct {
	logger          logger.Logger
	diseasesService Service
}

// NewHandler returns a new user Handler instance.
func NewHandler(logger logger.Logger, diseasesService Service) handler.Hand {
	return &Handler{
		logger:          logger,
		diseasesService: diseasesService,
	}
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, diseaseURL, h.GetDiseaseById)
	router.HandlerFunc(http.MethodPost, diseasesURL, h.CreateDisease)
	router.HandlerFunc(http.MethodPut, diseaseURL, h.UpdateDisease)
	router.HandlerFunc(http.MethodDelete, diseaseURL, h.DeleteDisease)
}

func (h *Handler) GetDiseaseById(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("GET DISEASE BY ID")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
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

	response.JSON(w, http.StatusOK, disease)
}

func (h *Handler) CreateDisease(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("CREATE DISEASE")

	var input CreateDiseaseDTO

	disease, err := h.diseasesService.Create(r.Context(), &input)
	if err != nil {
		response.InternalError(w, fmt.Sprintf("cannot create disease: %v", err), "")
		return
	}

	response.JSON(w, http.StatusCreated, disease)
}

func (h *Handler) UpdateDisease(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("UPDATE DISEASE")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	var input UpdateDiseaseDTO

	input.ID = id

	err = h.diseasesService.Update(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteDisease(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("DELETE DISEASE")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	err = h.diseasesService.Delete(r.Context(), id)
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
