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

/// Структура Handler представляющая собой обработчик объекта portfolioService для портфолио докторов \\\

type Handler struct {
	logger           logger.Logger
	portfolioService Service
}

/// Структура NewHandler возвращает новый экземпляр Handler инициализируя переданные в него аргументы \\\

func NewHandler(logger logger.Logger, portfolioService Service) handler.Hand {
	return &Handler{
		logger:           logger,
		portfolioService: portfolioService,
	}
}

/// Структура Register регистрирует новые запросы для портфолио докторов \\\

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, portfolioURL, h.GetPortfolioById)
	router.HandlerFunc(http.MethodPost, portfoliosURL, h.CreatePortfolio)
	router.HandlerFunc(http.MethodPut, portfolioURL, h.UpdatePortfolio)
	router.HandlerFunc(http.MethodDelete, portfolioURL, h.DeletePortfolio)
}

/// Функция GetPortfolioById получает портфолио по его id \\\

func (h *Handler) GetPortfolioById(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: GET PORTFOLIO BY ID")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	h.logger.Printf("Input: %+v\n", id)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	/// Вызов функции GetById передавая ей id портфолио \\\
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
	h.logger.Info("GOT PORTFOLIO BY ID")
	response.JSON(w, http.StatusOK, portfolio)
}

/// Функция CreatePortfolio создает портфолио по полученным данным из input \\\

func (h *Handler) CreatePortfolio(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: CREATE PORTFOLIO")

	var input CreatePortfolioDTO

	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.logger.Printf("Input: %+v\n", &input)

	/// Вызов функции Create передавая ей полученные значения и ссылку на структуру input \\\
	portfolio, err := h.portfolioService.Create(r.Context(), &input)
	if err != nil {
		response.InternalError(w, fmt.Sprintf("cannot create portfolio: %v", err), "")
		return
	}
	h.logger.Info("PORTFOLIO CREATED")
	response.JSON(w, http.StatusCreated, portfolio)
}

/// Функция UpdatePortfolio обновляет портфолио по его id и полученным данным из input \\\

func (h *Handler) UpdatePortfolio(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: UPDATE PORTFOLIO")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	var input UpdatePortfolioDTO

	/// Чтение JSON данных из тела входящего запроса r и декодирование их в переменную input \\\
	input.ID = id
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.logger.Printf("Input: %+v\n", &input)

	/// Вызов функции Update передавая ей полученные значения и ссылку на структуру input \\\
	err = h.portfolioService.Update(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
	}
	h.logger.Info("PORTFOLIO UPDATED")
	response.JSON(w, http.StatusOK, "PORTFOLIO UPDATED")
}

/// Функция DeletePortfolio удаляет портфолио по его id \\\

func (h *Handler) DeletePortfolio(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HANDLER: DELETE PORTFOLIO")

	/// Принимает объект r, представляющий HTTP-запрос, и извлекает параметр ID из URL \\\
	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	/// Вызов функции Delete передавая ей полученное значение id \\\
	err = h.portfolioService.Delete(id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		response.InternalError(w, err.Error(), "wrong on the server")
		return
	}
	h.logger.Info("PORTFOLIO DELETED")
	response.JSON(w, http.StatusOK, "PORTFOLIO DELETED")
}
