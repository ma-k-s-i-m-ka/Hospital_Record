package portfolio

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
	portfoliosURL = "/hospital_record/portfolios"
	portfolioURL  = "/hospital_record/portfolios/:id"
)

// Handler handles requests specified to user service.
type Handler struct {
	logger           logger.Logger
	portfolioService Service
}

// NewHandler returns a new user Handler instance.
func NewHandler(logger logger.Logger, portfolioService Service) handler.Hand {
	return &Handler{
		logger:           logger,
		portfolioService: portfolioService,
	}
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, portfolioURL, h.GetPortfolioById)
	router.HandlerFunc(http.MethodPost, portfoliosURL, h.CreatePortfolio)
	router.HandlerFunc(http.MethodPut, portfolioURL, h.UpdatePortfolio)
	router.HandlerFunc(http.MethodDelete, portfolioURL, h.DeletePortfolio)
}

func (h *Handler) GetPortfolioById(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("GET PORTFOLIO BY ID")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
	portfolio, err := h.portfolioService.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		h.logger.Error(err)
		response.InternalError(w, err.Error(), "")
		return
	}

	response.JSON(w, http.StatusOK, portfolio)
}

func (h *Handler) CreatePortfolio(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: CREATE PORTFOLIO")

	var input CreatePortfolioDTO
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.logger.Printf("Input: %+v\n", &input)
	portfolio, err := h.portfolioService.Create(r.Context(), &input)
	if err != nil {
		response.InternalError(w, fmt.Sprintf("cannot create portfolio: %v", err), "")
		return
	}
	response.JSON(w, http.StatusCreated, portfolio)
}

func (h *Handler) UpdatePortfolio(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: UPDATE PORTFOLIO")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	var input UpdatePortfolioDTO

	input.ID = id
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.logger.Printf("Input: %+v\n", &input)
	err = h.portfolioService.Update(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeletePortfolio(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: DELETE PORTFOLIO")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	err = h.portfolioService.Delete(id)
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
