package doctor

import (
	"HospitalRecord/app/internal/domain/apperror"
	"HospitalRecord/app/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"strings"
	"time"
)

var _ Storage = &DoctorStorage{}

type DoctorStorage struct {
	logger         logger.Logger
	conn           *pgx.Conn
	requestTimeout time.Duration
}

func NewStorage(storage *pgx.Conn, requestTimeout int) Storage {
	return &DoctorStorage{
		logger:         logger.GetLogger(),
		conn:           storage,
		requestTimeout: time.Duration(requestTimeout) * time.Second,
	}
}

func (d *DoctorStorage) Create(doctor *Doctor) (*Doctor, error) {

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	row := d.conn.QueryRow(ctx,
		`INSERT INTO doctors (name, surname, image_id, gender, rating, age,recording_is_available, specialization_id, portfolio_id)
			 VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9) 
			 RETURNING id`,
		doctor.Name, doctor.Surname, doctor.ImageID, doctor.Gender, doctor.Rating, doctor.Age, doctor.RecordingIsAvailable, doctor.SpecializationID, doctor.PortfolioID)

	err := row.Scan(&doctor.ID)
	if err != nil {
		err = fmt.Errorf("failed to execute create doctor query: %v", err)
		d.logger.Error(err)
		return nil, err
	}
	return doctor, nil
}

func (d *DoctorStorage) FindAll() ([]Doctor, error) {

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	rows, err := d.conn.Query(ctx,
		`SELECT * FROM doctor`)
	if err != nil {
		err = fmt.Errorf("failed to SELLECT: %v", err)
		d.logger.Error(err)
		return nil, err
	}

	doctors := make([]Doctor, 0)

	for rows.Next() {

		var doctor Doctor

		err = rows.Scan(
			&doctor.ID, &doctor.Name, &doctor.Surname, &doctor.Patronymic, &doctor.ImageID, &doctor.Gender,
			&doctor.Rating, &doctor.Age, &doctor.RecordingIsAvailable, &doctor.SpecializationID, &doctor.PortfolioID,
		)

		if err != nil {
			err = fmt.Errorf("failed to execute find all doctors query: %v", err)
			d.logger.Error(err)
			return nil, err
		}
		doctors = append(doctors, doctor)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return doctors, nil
}

func (d *DoctorStorage) FindAllAvailable(id int64, recordingIsAvailable bool) ([]Doctor, error) {

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	rows, err := d.conn.Query(ctx,
		`SELECT * FROM doctor
			 WHERE specialization_id=$1 AND recording_is_available=$2`,
		id, recordingIsAvailable)
	if err != nil {
		err = fmt.Errorf("failed to SELLECT: %v", err)
		d.logger.Error(err)
		return nil, err
	}

	doctors := make([]Doctor, 0)

	for rows.Next() {

		var doctor Doctor

		err = rows.Scan(
			&doctor.ID, &doctor.Name, &doctor.Surname, &doctor.Patronymic, &doctor.ImageID, &doctor.Gender,
			&doctor.Rating, &doctor.Age, &doctor.RecordingIsAvailable, &doctor.SpecializationID, &doctor.PortfolioID,
		)

		if err != nil {
			err = fmt.Errorf("failed to execute find all available doctors query: %v", err)
			d.logger.Error(err)
			return nil, err
		}
		doctors = append(doctors, doctor)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return doctors, nil
}

func (d *DoctorStorage) FindByPortfolioId(id int64) (*Doctor, error) {

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	row := d.conn.QueryRow(ctx,
		`SELECT * FROM doctor
			 WHERE portfolio_id = $1`, id)

	doctor := &Doctor{}

	err := row.Scan(
		&doctor.ID, &doctor.Name, &doctor.Surname, &doctor.Patronymic, &doctor.ImageID, &doctor.Gender,
		&doctor.Rating, &doctor.Age, &doctor.RecordingIsAvailable, &doctor.SpecializationID, &doctor.PortfolioID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrEmptyString
		}
		err = fmt.Errorf("failed to execute find user by portfolioid query: %v", err)
		d.logger.Error(err)
		return nil, err
	}

	return doctor, nil
}

func (d *DoctorStorage) FindById(id int64) (*Doctor, error) {

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	row := d.conn.QueryRow(ctx,
		`SELECT * FROM doctor
			 WHERE id = $1`, id)

	doctor := &Doctor{}

	err := row.Scan(
		&doctor.ID, &doctor.Name, &doctor.Surname, &doctor.Patronymic, &doctor.ImageID, &doctor.Gender,
		&doctor.Rating, &doctor.Age, &doctor.RecordingIsAvailable, &doctor.SpecializationID, &doctor.PortfolioID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrEmptyString
		}
		err = fmt.Errorf("failed to execute find doctor by id query: %v", err)
		d.logger.Error(err)
		return nil, err
	}

	return doctor, nil
}

func (d *DoctorStorage) Update(doctor *UpdateDoctorDTO) error {

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	result, err := d.conn.Exec(ctx,
		`UPDATE doctors
			SET patronymic=$1, image_id=$2, rating=$3, age=$4, recording_is_available=$5
			WHERE id =$6`,
		doctor.Patronymic, doctor.ImageID, doctor.Rating, doctor.Age, doctor.RecordingIsAvailable, doctor.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.ErrEmptyString
		}
		err = fmt.Errorf("failed to execute update doctor query: %v", err)
		d.logger.Error(err)
		return err
	}

	if result.RowsAffected() == 0 {
		return apperror.ErrEmptyString
	}
	return nil
}

func (d *DoctorStorage) PartiallyUpdate(doctor *PartiallyUpdateDoctorDTO) error {

	values := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if doctor.ImageID != nil {
		values = append(values, fmt.Sprintf("image_id=$%d", argId))
		args = append(args, *doctor.ImageID)
		argId++
	}

	if doctor.Rating != nil {
		values = append(values, fmt.Sprintf("rating=$%d", argId))
		args = append(args, *doctor.Rating)
		argId++
	}

	if doctor.RecordingIsAvailable != nil {
		values = append(values, fmt.Sprintf("record=$%d", argId))
		args = append(args, *doctor.RecordingIsAvailable)
		argId++
	}

	valuesQuery := strings.Join(values, ", ")
	query := fmt.Sprintf("UPDATE doctor  SET %s WHERE id = $%d", valuesQuery, argId)
	args = append(args, doctor.ID)

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	result, err := d.conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update doctor partially: %v", err)
	}

	if result.RowsAffected() == 0 {
		return apperror.ErrEmptyString
	}
	return nil
}

func (d *DoctorStorage) Delete(id int64) error {

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	result, err := d.conn.Exec(ctx,
		`DELETE FROM doctors WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete doctor: %v", err)
	}

	if result.RowsAffected() == 0 {
		return apperror.ErrEmptyString
	}
	return nil
}
