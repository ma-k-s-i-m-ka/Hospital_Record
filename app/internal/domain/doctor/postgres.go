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

/// Структура DoctorStorage содержащая поля для работы с БД \\\

type DoctorStorage struct {
	logger         logger.Logger
	conn           *pgx.Conn
	requestTimeout time.Duration
}

/// Структура NewStorage возвращает новый экземпляр DoctorStorage инициализируя переданные в него аргументы \\\

func NewStorage(storage *pgx.Conn, requestTimeout int) Storage {
	return &DoctorStorage{
		logger:         logger.GetLogger(),
		conn:           storage,
		requestTimeout: time.Duration(requestTimeout) * time.Second,
	}
}

/// Функция Create для сущности DoctorStorage создает записи докторов в БД \\\

func (d *DoctorStorage) Create(doctor *Doctor) (*Doctor, error) {
	d.logger.Info("POSTGRES: CREATE DOCTOR")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	row := d.conn.QueryRow(ctx,
		`INSERT INTO doctors (name, surname, image_id, gender, rating, age,recording_is_available, specialization_id, portfolio_id)
			 VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9) 
			 RETURNING id`,
		doctor.Name, doctor.Surname, doctor.ImageID, doctor.Gender, doctor.Rating, doctor.Age, doctor.RecordingIsAvailable, doctor.SpecializationID, doctor.PortfolioID)

	/// Сканирование полученных значений из БД \\\
	err := row.Scan(&doctor.ID)
	if err != nil {
		err = fmt.Errorf("failed to execute create doctor query: %v", err)
		d.logger.Error(err)
		return nil, err
	}
	return doctor, nil
}

/// Функция FindAll для сущности DoctorStorage находит всех докторов в БД \\\

func (d *DoctorStorage) FindAll() ([]Doctor, error) {
	d.logger.Info("POSTGRES: GET ALL DOCTORS")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	rows, err := d.conn.Query(ctx,
		`SELECT * FROM doctors`)
	if err != nil {
		err = fmt.Errorf("failed to SELLECT: %v", err)
		d.logger.Error(err)
		return nil, err
	}
	/// Создание пустого слайса для хранения всех хаписей \\\
	doctors := make([]Doctor, 0)

	/// Цикл создающий и записывающий новый экземпляр доктора \\\
	for rows.Next() {

		var doctor Doctor

		/// Сканирование полученных значений из БД \\\
		err = rows.Scan(
			&doctor.ID, &doctor.Name, &doctor.Surname, &doctor.Patronymic, &doctor.ImageID, &doctor.Gender,
			&doctor.Rating, &doctor.Age, &doctor.RecordingIsAvailable, &doctor.SpecializationID, &doctor.PortfolioID,
		)

		if err != nil {
			err = fmt.Errorf("failed to execute find all doctors query: %v", err)
			d.logger.Error(err)
			return nil, err
		}
		/// Добавление доктора в слайс \\\
		doctors = append(doctors, doctor)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return doctors, nil
}

/// Функция FindAllAvailable для сущности DoctorStorage находит всех свободных докторов по специализации в БД \\\

func (d *DoctorStorage) FindAllAvailable(id int64, recordingIsAvailable bool) ([]Doctor, error) {
	d.logger.Info("POSTGRES: GET ALL AVAILABLE DOCTORS")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	rows, err := d.conn.Query(ctx,
		`SELECT * FROM doctors
			 WHERE specialization_id=$1 AND recording_is_available=$2`,
		id, recordingIsAvailable)
	if err != nil {
		err = fmt.Errorf("failed to SELLECT: %v", err)
		d.logger.Error(err)
		return nil, err
	}

	/// Создание пустого слайса для хранения всех хаписей \\\
	doctors := make([]Doctor, 0)

	/// Цикл создающий и записывающий новый экземпляр доктора \\\
	for rows.Next() {

		var doctor Doctor

		/// Сканирование полученных значений из БД \\\
		err = rows.Scan(
			&doctor.ID, &doctor.Name, &doctor.Surname, &doctor.Patronymic, &doctor.ImageID, &doctor.Gender,
			&doctor.Rating, &doctor.Age, &doctor.RecordingIsAvailable, &doctor.SpecializationID, &doctor.PortfolioID,
		)

		if err != nil {
			err = fmt.Errorf("failed to execute find all available doctors query: %v", err)
			d.logger.Error(err)
			return nil, err
		}
		/// Добавление доктора в слайс \\\
		doctors = append(doctors, doctor)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return doctors, nil
}

/// Функция FindByPortfolioId для сущности DoctorStorage получает записи доктора из БД по id портфолио доктора \\\

func (d *DoctorStorage) FindByPortfolioId(id int64) (*Doctor, error) {
	d.logger.Info("POSTGRES: GET DOCTOR BY PORTFOLIO ID")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	row := d.conn.QueryRow(ctx,
		`SELECT * FROM doctors
			 WHERE portfolio_id = $1`, id)

	doctor := &Doctor{}

	/// Сканирование полученных значений из БД \\\
	err := row.Scan(
		&doctor.ID, &doctor.Name, &doctor.Surname, &doctor.Patronymic, &doctor.ImageID, &doctor.Gender,
		&doctor.Rating, &doctor.Age, &doctor.RecordingIsAvailable, &doctor.SpecializationID, &doctor.PortfolioID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		err = fmt.Errorf("failed to execute find user by portfolio id query: %v", err)
		d.logger.Error(err)
		return nil, err
	}

	return doctor, nil
}

/// Функция FindById для сущности DoctorStorage получает записи доктора из БД по id доктора \\\

func (d *DoctorStorage) FindById(id int64) (*Doctor, error) {
	d.logger.Info("POSTGRES: GET DOCTOR BY ID")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()
	d.logger.Printf("Input: %+v\n", id)

	/// Выполнение запроса к БД \\\
	row := d.conn.QueryRow(ctx,
		`SELECT * FROM doctors
			 WHERE id = $1`, id)

	doctor := &Doctor{}

	/// Сканирование полученных значений из БД \\\
	err := row.Scan(
		&doctor.ID, &doctor.Name, &doctor.Surname, &doctor.Patronymic,
		&doctor.ImageID, &doctor.Gender, &doctor.Rating, &doctor.Age,
		&doctor.RecordingIsAvailable, &doctor.SpecializationID,
		&doctor.PortfolioID,
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

/// Функция Update для сущности DoctorStorage обновляет записи о докторе в БД \\\

func (d *DoctorStorage) Update(doctor *UpdateDoctorDTO) error {
	d.logger.Info("POSTGRES: UPDATE DOCTOR")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
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

/// Функция PartiallyUpdate для сущности DoctorStorage частично обновляет записи о докторе в БД \\\

func (d *DoctorStorage) PartiallyUpdate(doctor *PartiallyUpdateDoctorDTO) error {
	d.logger.Info("POSTGRES: PARTIALLY UPDATE DOCTOR")

	/// Создание пустого слайса для хранения обновляемых строк \\\
	values := make([]string, 0)

	/// Создание пустого слайса для хранения аргументов запроса \\\
	args := make([]interface{}, 0)
	argId := 1

	/// Проверки на наличие новых значений \\\
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
		values = append(values, fmt.Sprintf("recording_is_available=$%d", argId))
		args = append(args, *doctor.RecordingIsAvailable)
		argId++
	}

	/// Формирование строки со всеми измененными полями и их значениями \\\
	valuesQuery := strings.Join(values, ", ")
	query := fmt.Sprintf("UPDATE doctors  SET %s WHERE id = $%d", valuesQuery, argId)
	args = append(args, doctor.ID)

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	result, err := d.conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update doctor partially: %v", err)
	}

	if result.RowsAffected() == 0 {
		return apperror.ErrEmptyString
	}
	return nil
}

/// Функция Delete для сущности DoctorStorage удаляет записи о докторое из БД \\\

func (d *DoctorStorage) Delete(id int64) error {
	d.logger.Info("POSTGRES: DELETE DOCTOR")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
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
