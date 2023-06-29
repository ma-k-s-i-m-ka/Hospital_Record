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

type DiseaseStorage struct {
	logger         logger.Logger
	conn           *pgx.Conn
	requestTimeout time.Duration
}

func NewStorage(storage *pgx.Conn, requestTimeout int) Storage {
	return &DiseaseStorage{
		logger:         logger.GetLogger(),
		conn:           storage,
		requestTimeout: time.Duration(requestTimeout) * time.Second,
	}
}

func (d *DiseaseStorage) Create(disease *Disease) (*Disease, error) {

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	row := d.conn.QueryRow(ctx,
		`INSERT INTO disease (body_part, description)
			 VALUES($1,$2) 
			 RETURNING id`,
		disease.BodyPart, disease.Description)

	err := row.Scan(&disease.ID)
	if err != nil {
		err = fmt.Errorf("failed to execute create disease query: %v", err)
		d.logger.Error(err)
		return nil, err
	}
	return disease, nil
}

func (d *DiseaseStorage) FindById(id int64) (*Disease, error) {

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	row := d.conn.QueryRow(ctx,
		`SELECT * FROM disease
			 WHERE id = $1`, id)

	disease := &Disease{}

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

func (d *DiseaseStorage) Update(disease *UpdateDiseaseDTO) error {

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

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

func (d *DiseaseStorage) Delete(id int64) error {

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

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
