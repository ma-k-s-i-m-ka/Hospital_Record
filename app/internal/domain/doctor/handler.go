package doctor

import (
	"HospitalRecord/app/internal/domain/apperror"
	"HospitalRecord/app/internal/domain/handler"
	"HospitalRecord/app/internal/domain/response"
	"HospitalRecord/app/pkg/logger"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
)

const (
	doctorsURL             = "/hospital_record/doctors"
	doctorsAllURL          = "/hospital_record/all_doctors"
	doctorsAvailableURL    = "/hospital_record/doctors/available/:id"
	doctorURL              = "/hospital_record/doctors/profile/:id"
	doctorByPortfolioIdURL = "/hospital_record/doctors/portfolio/:id"
	doctorImageURL         = "/hospital_record/doctor/image"
)

/// Структура Handler представляющая собой обработчик объекта doctorService для докторов \\\

type Handler struct {
	logger        logger.Logger
	doctorService Service
}

/// Структура NewHandler возвращает новый экземпляр Handler инициализируя переданные в него аргументы \\\

func NewHandler(logger logger.Logger, doctorService Service) handler.Hand {
	return &Handler{
		logger:        logger,
		doctorService: doctorService,
	}
}

/// Структура Register регистрирует новые запросы для докторов \\\

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, doctorURL, h.GetDoctorById)
	router.HandlerFunc(http.MethodGet, doctorByPortfolioIdURL, h.GetDoctorByPortfolioId)
	router.HandlerFunc(http.MethodGet, doctorsAllURL, h.FindAllDoctors)
	router.HandlerFunc(http.MethodGet, doctorsAvailableURL, h.FindAllAvailableDoctors)
	router.HandlerFunc(http.MethodPost, doctorsURL, h.CreateDoctor)
	router.HandlerFunc(http.MethodPost, doctorImageURL, h.CreateImage)
	router.HandlerFunc(http.MethodPut, doctorURL, h.UpdateDoctor)
	router.HandlerFunc(http.MethodPatch, doctorURL, h.PartiallyUpdateDoctor)
	router.HandlerFunc(http.MethodDelete, doctorURL, h.DeleteDoctor)
}

/// Функция GetDoctorById получает доктора по его id \\\

func (h *Handler) GetDoctorById(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: GET DOCTOR BY ID")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	h.logger.Printf("Input: %+v\n", id)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	/// Вызов функции GetById передавая ей id доктора \\\
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
	h.logger.Info("GOT DOCTOR BY ID")
	response.JSON(w, http.StatusOK, doctor)
}

/// Функция GetDoctorByPortfolioId получает доктора по его id портфолио \\\

func (h *Handler) GetDoctorByPortfolioId(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: GET DOCTOR BY PORTFOLIO ID")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	h.logger.Printf("Input: %+v\n", id)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	/// Вызов функции GetByPortfolioId передавая ей id портфолио \\\
	doctor, err := h.doctorService.GetByPortfolioId(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperror.ErrRepeatedPortfolioId) {
			response.NotFound(w)
			return
		}
		response.BadRequest(w, err.Error(), "")
		return
	}
	h.logger.Info("GOT DOCTOR BY PORTFOLIO ID")
	response.JSON(w, http.StatusOK, doctor)
}

/// Функция FindAllDoctors получает всех докторов \\\

func (h *Handler) FindAllDoctors(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: GET ALL DOCTORS")

	/// Вызов функции FindAll \\\
	doctors, err := h.doctorService.FindAll()
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
	h.logger.Info("GOT ALL DOCTORS")
	response.JSON(w, http.StatusOK, doctors)
}

/// Функция FindAllAvailableDoctors получает всех свободных докторов по id их специализации \\\

func (h *Handler) FindAllAvailableDoctors(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: GET ALL AVAILABLE DOCTORS")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	/// Вызов функции FindAllAvailable передавая ей id нужной специализации \\\
	doctors, err := h.doctorService.FindAllAvailable(r.Context(), id, true)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
	h.logger.Info("GOT ALL AVAILABLE DOCTORS")
	response.JSON(w, http.StatusOK, doctors)
}

/// Функция CreateImage создает изображение доктора в папке doctorimages \\\

func (h *Handler) CreateImage(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: CREATE DOCTOR IMAGE")

	/// Принимает объект r типа form-data \\\
	file, header, err := r.FormFile("image")
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
	defer file.Close()

	/// Создаем путь и новый файл для сохранения изображений \\\
	imagePath := "./doctorimages/" + header.Filename
	out, err := os.Create(imagePath)
	if err != nil {
		response.InternalError(w, fmt.Sprintf("error saving image: %v", err), "")
		return
	}
	defer out.Close()

	/// Копируем загруженный файл в новый созданный файл на севрере по созданному путю \\\
	_, err = io.Copy(out, file)
	if err != nil {
		response.InternalError(w, fmt.Sprintf("error coping image: %v", err), "")
		return
	}
	h.logger.Info("DOCTOR IMAGE CREATED")
	response.JSON(w, http.StatusCreated, "IMAGE CREATED")
}

/// Функция CreateDoctor создает доктора по полученным данным из input \\\

func (h *Handler) CreateDoctor(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: CREATE DOCTOR")

	var input CreateDoctorDTO

	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.logger.Printf("Input: %+v\n", &input)

	/// Вызов функции Create передавая ей полученные значения и ссылку на структуру input \\\
	doctor, err := h.doctorService.Create(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrRepeatedPortfolioId) {
			response.BadRequest(w, err.Error(), "")
			return
		}
		response.InternalError(w, fmt.Sprintf("cannot create doctor: %v", err), "")
		return
	}
	h.logger.Info("DOCTOR CREATED")
	response.JSON(w, http.StatusCreated, doctor)
}

/// Функция UpdateDoctor обновляет доктора по его id и полученным данным из input \\\

func (h *Handler) UpdateDoctor(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: UPDATE DOCTOR")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	var input UpdateDoctorDTO

	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	input.ID = id
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.logger.Printf("Input: %+v\n", &input)

	/// Вызов функции Update передавая ей полученные значения и ссылку на структуру input \\\
	err = h.doctorService.Update(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
	}
	h.logger.Info("DOCTOR UPDATED")
	response.JSON(w, http.StatusOK, "DOCTOR UPDATED")
}

/// Функция PartiallyUpdateDoctor частично обновляет доктора по его id и полученным данным из input \\\

func (h *Handler) PartiallyUpdateDoctor(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: PARTIALLY UPDATE DOCTOR")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
	var input PartiallyUpdateDoctorDTO
	input.ID = id

	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	/// Вызов функции PartiallyUpdate передавая ей полученные значения и ссылку на структуру input \\\
	err = h.doctorService.PartiallyUpdate(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
	}
	h.logger.Info("DOCTOR PARTIALLY UPDATED")
	response.JSON(w, http.StatusOK, "DOCTOR PARTIALLY UPDATED")
}

/// Функция DeleteDoctor удаляет доктора по его id \\\

func (h *Handler) DeleteDoctor(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: DELETE DOCTOR")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	/// Вызов функции Delete передавая ей полученное значение id \\\
	err = h.doctorService.Delete(id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		response.InternalError(w, err.Error(), "wrong on the server")
		return
	}
	h.logger.Info("DOCTOR DELETED")
	response.JSON(w, http.StatusOK, "DOCTOR DELETED")
}
