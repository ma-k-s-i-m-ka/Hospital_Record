package disease

import (
	"HospitalRecord/app/internal/domain/apperror"
	"HospitalRecord/app/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"time"
)

var _ Storage = &DiseaseStorage{}

/// Структура DiseaseStorage содержащая поля для работы с БД \\\

type DiseaseStorage struct {
	logger         logger.Logger
	conn           *pgx.Conn
	requestTimeout time.Duration
}

/// Структура NewStorage возвращает новый экземпляр DiseaseStorage инициализируя переданные в него аргументы \\\

func NewStorage(storage *pgx.Conn, requestTimeout int) Storage {
	return &DiseaseStorage{
		logger:         logger.GetLogger(),
		conn:           storage,
		requestTimeout: time.Duration(requestTimeout) * time.Second,
	}
}

/// Функция Create для сущности DiseaseStorage создает записи о болезни в БД \\\

func (d *DiseaseStorage) Create(disease *Disease) (*Disease, error) {
	d.logger.Info("POSTGRES: CREATE DISEASE")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	row := d.conn.QueryRow(ctx,
		`INSERT INTO disease (body_part, description)
			 VALUES($1,$2) 
			 RETURNING id`,
		disease.BodyPart, disease.Description)

	/// Сканирование полученных значений из БД \\\
	err := row.Scan(&disease.ID)
	if err != nil {
		err = fmt.Errorf("failed to execute create disease query: %v", err)
		d.logger.Error(err)
		return nil, err
	}
	return disease, nil
}

/// Функция FindById для сущности DiseaseStorage получает записи о болезни из БД по id болезни\\\

func (d *DiseaseStorage) FindById(id int64) (*Disease, error) {
	d.logger.Info("POSTGRES: GET DISEASE BY ID")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	row := d.conn.QueryRow(ctx,
		`SELECT * FROM disease
			 WHERE id = $1`, id)

	disease := &Disease{}

	/// Сканирование полученных значений из БД \\\
	err := row.Scan(
		&disease.ID, &disease.BodyPart, &disease.Description,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrEmptyString
		}
		err = fmt.Errorf("failed to execute find disease by id query: %v", err)
		d.logger.Error(err)
		return nil, err
	}

	return disease, nil
}

/// Функция Update для сущности DiseaseStorage обновляет записи о болезни в БД \\\

func (d *DiseaseStorage) Update(disease *UpdateDiseaseDTO) error {
	d.logger.Info("POSTGRES: UPDATE DISEASE")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	result, err := d.conn.Exec(ctx,
		`UPDATE disease
			SET body_part=$1, description=$2
			WHERE id =$3`,
		disease.BodyPart, disease.Description, disease.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.ErrEmptyString
		}
		err = fmt.Errorf("failed to execute update disease query: %v", err)
		d.logger.Error(err)
		return err
	}

	if result.RowsAffected() == 0 {
		return apperror.ErrEmptyString
	}
	return nil
}

/// Функция Delete для сущности DiseaseStorage удаляет записи о болезни из БД \\\

func (d *DiseaseStorage) Delete(id int64) error {
	d.logger.Info("POSTGRES: DELETE DISEASE")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	result, err := d.conn.Exec(ctx,
		`DELETE FROM disease WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete disease: %v", err)
	}

	if result.RowsAffected() == 0 {
		return apperror.ErrEmptyString
	}
	return nil
}
