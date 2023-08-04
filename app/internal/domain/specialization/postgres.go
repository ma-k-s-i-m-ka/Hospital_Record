package specialization

import (
	"HospitalRecord/app/internal/domain/apperror"
	"HospitalRecord/app/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"time"
)

var _ Storage = &SpecializationStorage{}

/// Структура SpecializationStorage содержащая поля для работы с БД \\\

type SpecializationStorage struct {
	logger         logger.Logger
	conn           *pgx.Conn
	requestTimeout time.Duration
}

/// Структура NewStorage возвращает новый экземпляр SpecializationStorage инициализируя переданные в него аргументы \\\

func NewStorage(storage *pgx.Conn, requestTimeout int) Storage {
	return &SpecializationStorage{
		logger:         logger.GetLogger(),
		conn:           storage,
		requestTimeout: time.Duration(requestTimeout) * time.Second,
	}
}

/// Функция Create для сущности SpecializationStorage создает записи специализации в БД \\\

func (d *SpecializationStorage) Create(specialization *Specialization) (*Specialization, error) {
	d.logger.Info("POSTGRES: CREATE SPECIALIZATION")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	row := d.conn.QueryRow(ctx,
		`INSERT INTO specialization (name_specialization)
			 VALUES($1) 
			 RETURNING id`,
		specialization.Name)

	/// Сканирование полученных значений из БД \\\
	err := row.Scan(&specialization.ID)
	if err != nil {
		err = fmt.Errorf("failed to execute create specialization query: %v", err)
		d.logger.Error(err)
		return nil, err
	}
	return specialization, nil
}

/// Функция FindById для сущности SpecializationStorage получает записи специализации из БД по id \\\

func (d *SpecializationStorage) FindById(id int64) (*Specialization, error) {
	d.logger.Info("POSTGRES: GET SPECIALIZATION BY ID")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	row := d.conn.QueryRow(ctx,
		`SELECT * FROM specialization
			 WHERE id = $1`, id)

	specialization := &Specialization{}

	/// Сканирование полученных значений из БД \\\
	err := row.Scan(
		&specialization.ID, &specialization.Name,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrEmptyString
		}
		err = fmt.Errorf("failed to execute find specialization by id query: %v", err)
		d.logger.Error(err)
		return nil, err
	}
	return specialization, nil
}

/// Функция Update для сущности SpecializationStorage обновляет записи о специализации в БД \\\

func (d *SpecializationStorage) Update(specialization *UpdateSpecializationDTO) error {
	d.logger.Info("POSTGRES: UPDATE SPECIALIZATION")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	result, err := d.conn.Exec(ctx,
		`UPDATE specialization
			SET name_specialization=$1
			WHERE id =$2`,
		specialization.Name, specialization.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.ErrEmptyString
		}
		err = fmt.Errorf("failed to execute update specialization query: %v", err)
		d.logger.Error(err)
		return err
	}
	if result.RowsAffected() == 0 {
		return apperror.ErrEmptyString
	}
	return nil
}

/// Функция Delete для сущности SpecializationStorage удаляет записи о специализации из БД \\\

func (d *SpecializationStorage) Delete(id int64) error {
	d.logger.Info("POSTGRES: DELETE SPECIALIZATION")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	result, err := d.conn.Exec(ctx,
		`DELETE FROM specialization WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete specialization: %v", err)
	}

	if result.RowsAffected() == 0 {
		return apperror.ErrEmptyString
	}
	return nil
}
