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

type SpecializationStorage struct {
	logger         logger.Logger
	conn           *pgx.Conn
	requestTimeout time.Duration
}

func NewStorage(storage *pgx.Conn, requestTimeout int) Storage {
	return &SpecializationStorage{
		logger:         logger.GetLogger(),
		conn:           storage,
		requestTimeout: time.Duration(requestTimeout) * time.Second,
	}
}

func (d *SpecializationStorage) Create(specialization *Specialization) (*Specialization, error) {

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	row := d.conn.QueryRow(ctx,
		`INSERT INTO specialization (name_specialization)
			 VALUES($1) 
			 RETURNING id`,
		specialization.Name)

	err := row.Scan(&specialization.ID)
	if err != nil {
		err = fmt.Errorf("failed to execute create specialization query: %v", err)
		d.logger.Error(err)
		return nil, err
	}
	return specialization, nil
}

func (d *SpecializationStorage) FindById(id int64) (*Specialization, error) {

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	row := d.conn.QueryRow(ctx,
		`SELECT * FROM specialization
			 WHERE id = $1`, id)

	specialization := &Specialization{}

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

func (d *SpecializationStorage) Update(specialization *UpdateSpecializationDTO) error {

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

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

func (d *SpecializationStorage) Delete(id int64) error {

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

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
