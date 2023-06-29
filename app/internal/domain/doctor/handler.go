package doctor

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
	doctorsURL = "/hospitalrecord/doctors"
	doctorURL  = "/hospitalrecord/doctors/:id"
)

// Handler handles requests specified to user service.
type Handler struct {
	logger        logger.Logger
	doctorService Service
}

func NewHandler(logger logger.Logger, doctorService Service) handler.Hand {
	return &Handler{
		logger:        logger,
		doctorService: doctorService,
	}
}

// Register registers new routes for router.
func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, doctorURL, h.GetDoctorById)
	router.HandlerFunc(http.MethodGet, doctorURL, h.GetDoctorByPortfolioId)
	router.HandlerFunc(http.MethodGet, doctorsURL, h.FindAllDoctors)
	router.HandlerFunc(http.MethodGet, doctorsURL, h.FindAllAvailableDoctors)
	router.HandlerFunc(http.MethodPost, doctorURL, h.CreateDoctor)
	router.HandlerFunc(http.MethodPut, doctorURL, h.UpdateDoctor)
	router.HandlerFunc(http.MethodPatch, doctorURL, h.PartiallyUpdateDoctor)
	router.HandlerFunc(http.MethodDelete, doctorURL, h.DeleteDoctor)
}

func (h *Handler) GetDoctorById(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("GET DOCTOR BY ID")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
	doctor, err := h.doctorService.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		h.logger.Error(err)
		response.InternalError(w, err.Error(), "")
		return
	}

	response.JSON(w, http.StatusOK, doctor)
}

func (h *Handler) GetDoctorByPortfolioId(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("GET DOCTOR BY PORTFOLIO ID")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	doctor, err := h.doctorService.GetByPortfolioId(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperror.ErrRepeatedPortfolioId) {
			response.NotFound(w)
			return
		}
		response.BadRequest(w, err.Error(), "")
		return
	}
	response.JSON(w, http.StatusOK, doctor)
}

func (h *Handler) FindAllDoctors(w http.ResponseWriter, r *http.Request) {
	doctors, err := h.doctorService.FindAll()
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
	response.JSON(w, http.StatusOK, doctors)
}

func (h *Handler) FindAllAvailableDoctors(w http.ResponseWriter, r *http.Request) {
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	doctors, err := h.doctorService.FindAllAvailable(r.Context(), id, true)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
	response.JSON(w, http.StatusOK, doctors)
}

func (h *Handler) CreateDoctor(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("CREATE DOCTOR")

	var input CreateDoctorDTO

	user, err := h.doctorService.Create(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrRepeatedPortfolioId) {
			response.BadRequest(w, err.Error(), "")
			return
		}
		response.InternalError(w, fmt.Sprintf("cannot create doctor: %v", err), "")
		return
	}

	response.JSON(w, http.StatusCreated, user)
}

func (h *Handler) UpdateDoctor(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("UPDATE DOCTOR")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	var input UpdateDoctorDTO

	input.ID = id

	err = h.doctorService.Update(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) PartiallyUpdateDoctor(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("PARTIALLY UPDATE DOCTOR")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	var input PartiallyUpdateDoctorDTO

	input.ID = id

	err = h.doctorService.PartiallyUpdate(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteDoctor(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("DELETE DOCTOR")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	err = h.doctorService.Delete(r.Context(), id)
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
