package portfolio

import (
	"HospitalRecord/app/internal/domain/apperror"
	"HospitalRecord/app/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"time"
)

var _ Storage = &PortfolioStorage{}

/// Структура PortfolioStorage содержащая поля для работы с БД \\\

type PortfolioStorage struct {
	logger         logger.Logger
	conn           *pgx.Conn
	requestTimeout time.Duration
}

/// Структура NewStorage возвращает новый экземпляр PortfolioStorage инициализируя переданные в него аргументы \\\

func NewStorage(storage *pgx.Conn, requestTimeout int) Storage {
	return &PortfolioStorage{
		logger:         logger.GetLogger(),
		conn:           storage,
		requestTimeout: time.Duration(requestTimeout) * time.Second,
	}
}

/// Функция Create для сущности PortfolioStorage создает записи портфолио в БД \\\

func (d *PortfolioStorage) Create(portfolio *Portfolio) (*Portfolio, error) {
	d.logger.Info("POSTGRES: CREATE PORTFOLIO")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	row := d.conn.QueryRow(ctx,
		`INSERT INTO portfolio (education, awards, work_experience)
			 VALUES($1,$2,$3) 
			 RETURNING id`,
		portfolio.Education, portfolio.Awards, portfolio.WorkExperience)

	/// Сканирование полученных значений из БД \\\
	err := row.Scan(&portfolio.ID)
	if err != nil {
		err = fmt.Errorf("failed to execute create portfolio query: %v", err)
		d.logger.Error(err)
		return nil, err
	}
	return portfolio, nil
}

/// Функция FindById для сущности PortfolioStorage получает записи портфолио из БД по id \\\

func (d *PortfolioStorage) FindById(id int64) (*Portfolio, error) {
	d.logger.Info("POSTGRES: GET PORTFOLIO BY ID")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	row := d.conn.QueryRow(ctx,
		`SELECT * FROM portfolio
			 WHERE id = $1`, id)

	portfolio := &Portfolio{}

	/// Сканирование полученных значений из БД \\\
	err := row.Scan(
		&portfolio.ID, &portfolio.Education, &portfolio.Awards, &portfolio.WorkExperience,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrEmptyString
		}
		err = fmt.Errorf("failed to execute find portfolio by id query: %v", err)
		d.logger.Error(err)
		return nil, err
	}
	return portfolio, nil
}

/// Функция Update для сущности PortfolioStorage обновляет записи о портфолио в БД \\\

func (d *PortfolioStorage) Update(portfolio *UpdatePortfolioDTO) error {
	d.logger.Info("POSTGRES: UPDATE PORTFOLIO")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	result, err := d.conn.Exec(ctx,
		`UPDATE portfolio
			SET education=$1, awards=$2, work_experience=$3
			WHERE id =$4`,
		portfolio.Education, portfolio.Awards, portfolio.WorkExperience, portfolio.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.ErrEmptyString
		}
		err = fmt.Errorf("failed to execute update portfolio query: %v", err)
		d.logger.Error(err)
		return err
	}
	if result.RowsAffected() == 0 {
		return apperror.ErrEmptyString
	}
	return nil
}

/// Функция Delete для сущности PortfolioStorage удаляет записи портфолио из БД \\\

func (d *PortfolioStorage) Delete(id int64) error {
	d.logger.Info("POSTGRES: DELETE PORTFOLIO")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	result, err := d.conn.Exec(ctx,
		`DELETE FROM portfolio WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete portfolio: %v", err)
	}

	if result.RowsAffected() == 0 {
		return apperror.ErrEmptyString
	}
	return nil
}
