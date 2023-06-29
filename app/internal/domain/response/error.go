package response

import (
	"HospitalRecord/app/internal/domain/apperror"
	"net/http"
)

func Error(w http.ResponseWriter, code int, message, developerMessage string) {
	appError := apperror.NewAppError(code, message, developerMessage)
	JSON(w, code, appError)
}
func BadRequest(w http.ResponseWriter, message, developerMessage string) {
	Error(w, http.StatusBadRequest, message, developerMessage)
}

func NotFound(w http.ResponseWriter) {
	JSON(w, http.StatusNotFound, apperror.ErrNotFound)
}

func InternalError(w http.ResponseWriter, message, developerMessage string) {
	Error(w, http.StatusInternalServerError, message, developerMessage)
}
