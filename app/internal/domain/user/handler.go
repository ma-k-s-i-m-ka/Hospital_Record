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

// Handler handles requests specified to user service.
type Handler struct {
	logger      logger.Logger
	userService Service
}

// NewHandler returns a new user Handler instance.
func NewHandler(logger logger.Logger, userService Service) handler.Hand {
	return &Handler{
		logger:      logger,
		userService: userService,
	}
}

// Register registers new routes for router.
func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, userByEmailURL, h.GetUserByEmail)
	router.HandlerFunc(http.MethodGet, userByPolicyNumber, h.GetUserByPolicyNumber)
	router.HandlerFunc(http.MethodPost, usersURL, h.CreateUser)
	router.HandlerFunc(http.MethodPut, userURL, h.UpdateUser)
	router.HandlerFunc(http.MethodPatch, userURL, h.PartiallyUpdateUser)
	router.HandlerFunc(http.MethodDelete, userURL, h.DeleteUser)
	router.HandlerFunc(http.MethodGet, userURL, h.GetUserById)
}

func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: GET USER BY ID")

	id, err := handler.ReadIdParam64(r)
	h.logger.Printf("Input: %+v\n", id)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
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
	response.JSON(w, http.StatusOK, user)
}

func (h *Handler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: GET USER BY EMAIL")

	email := r.URL.Query().Get("email")

	h.logger.Printf("Input: %+v\n", email)

	if email == "" {
		response.BadRequest(w, "empty email", "")
		return
	}

	user, err := h.userService.GetByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		response.BadRequest(w, err.Error(), "")
		return
	}
	response.JSON(w, http.StatusOK, user)
}

func (h *Handler) GetUserByPolicyNumber(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: GET USER BY POLICY NUMBER")

	policy := r.URL.Query().Get("policy_number")

	if policy == "" {
		response.BadRequest(w, "empty policy", "")
		return
	}

	user, err := h.userService.GetByPolicyNumber(r.Context(), policy)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		response.BadRequest(w, err.Error(), "")
		return
	}
	response.JSON(w, http.StatusOK, user)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: CREATE USER")

	var input CreateUserDTO
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.logger.Printf("Input: %+v\n", &input)
	user, err := h.userService.Create(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrRepeatedEmail) {
			response.BadRequest(w, err.Error(), "")
			return
		}
		response.InternalError(w, fmt.Sprintf("cannot create user: %v", err), "")
		return
	}

	response.JSON(w, http.StatusCreated, user)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: UPDATE USER")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
	var input UpdateUserDTO
	input.ID = id
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.logger.Printf("Input: %+v\n", &input)
	err = h.userService.Update(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: PARTIALLY UPDATE USER")
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
	var input PartiallyUpdateUserDTO
	input.ID = id
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	err = h.userService.PartiallyUpdate(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: DELETE USER")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	err = h.userService.Delete(id)
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
