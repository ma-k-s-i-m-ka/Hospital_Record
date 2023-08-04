package record

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

var _ Storage = &RecordStorage{}

/// Структура RecordStorage содержащая поля для работы с БД \\\

type RecordStorage struct {
	logger         logger.Logger
	conn           *pgx.Conn
	requestTimeout time.Duration
}

/// Структура NewStorage возвращает новый экземпляр RecordStorage инициализируя переданные в него аргументы \\\

func NewStorage(storage *pgx.Conn, requestTimeout int) Storage {
	return &RecordStorage{
		logger:         logger.GetLogger(),
		conn:           storage,
		requestTimeout: time.Duration(requestTimeout) * time.Second,
	}
}

/// Функция CreateRecord для сущности RecordStorage создает записи на прием в БД \\\

func (r *RecordStorage) CreateRecord(record *Record) (*Record, error) {
	r.logger.Info("POSTGRES: CREATE RECORD")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), r.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	row := r.conn.QueryRow(ctx,
		`INSERT INTO record (hospital_address, doctor_office, tagging, patients_id, doctor_id, specialization_id, time_record)
			 VALUES($1,$2,$3,$4,$5,$6,$7) 
			 RETURNING id`,
		record.HospitalAddress, record.DoctorOffice, record.Tagging, record.PatientsID, record.DoctorID, record.SpecializationID, record.TimeRecord)

	/// Сканирование полученных значений из БД \\\
	err := row.Scan(&record.ID)
	if err != nil {
		err = fmt.Errorf("failed to execute create record query: %v", err)
		r.logger.Error(err)
		return nil, err
	}
	return record, nil
}

/// Функция FindRecordByPatientsId для сущности RecordStorage получает записи на прием из БД по id пациента \\\

func (r *RecordStorage) FindRecordByPatientsId(id int64) (*Record, error) {
	r.logger.Info("POSTGRES: GET RECORD BY PATIENTS ID")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), r.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	row := r.conn.QueryRow(ctx,
		`SELECT * FROM record
			 WHERE patients_id = $1`, id)

	record := &Record{}

	/// Сканирование полученных значений из БД \\\
	err := row.Scan(&record.ID, &record.HospitalAddress, &record.DoctorOffice, &record.Tagging,
		&record.PatientsID, &record.DoctorID, &record.SpecializationID,
		&record.TimeRecord,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		err = fmt.Errorf("failed to execute find record by patients id query: %v", err)
		r.logger.Error(err)
		return nil, err
	}
	return record, nil
}

/// Функция FindRecordById для сущности RecordStorage получает записи на прием из БД по id \\\

func (r *RecordStorage) FindRecordById(id int64) (*Record, error) {
	r.logger.Info("POSTGRES: GET RECORD BY ID")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), r.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	row := r.conn.QueryRow(ctx,
		`SELECT * FROM record
			 WHERE id = $1`, id)

	record := &Record{}

	/// Сканирование полученных значений из БД \\\
	err := row.Scan(&record.ID, &record.HospitalAddress, &record.DoctorOffice, &record.Tagging,
		&record.PatientsID, &record.DoctorID, &record.SpecializationID,
		&record.TimeRecord,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrEmptyString
		}
		err = fmt.Errorf("failed to execute find record by id query: %v", err)
		r.logger.Error(err)
		return nil, err
	}

	return record, nil
}

/// Функция UpdateRecord для сущности RecordStorage обновляет записи на прием в БД \\\

func (r *RecordStorage) UpdateRecord(record *UpdateRecordDTO) error {
	r.logger.Info("POSTGRES: UPDATE RECORD")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), r.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	result, err := r.conn.Exec(ctx,
		`UPDATE record
			 SET hospital_address=$1, doctor_office=$2, tagging=$3, patients_id=$4, doctor_id=$5, specialization_id=$6, time_record=$7
			 WHERE id =$8`,
		record.HospitalAddress, record.DoctorOffice, record.Tagging, record.PatientsID, record.DoctorID, record.SpecializationID, record.TimeRecord, &record.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.ErrEmptyString
		}
		err = fmt.Errorf("failed to execute update record query: %v", err)
		r.logger.Error(err)
		return err
	}
	if result.RowsAffected() == 0 {
		return apperror.ErrEmptyString
	}
	return nil
}

/// Функция PartiallyUpdateRecord для сущности RecordStorage частично обновляет записи на прием в БД \\\

func (r *RecordStorage) PartiallyUpdateRecord(record *PartiallyUpdateRecordDTO) error {
	r.logger.Info("POSTGRES: PARTIALLY UPDATE RECORD")

	/// Создание пустого слайса для хранения обновляемых строк \\\
	values := make([]string, 0)

	/// Создание пустого слайса для хранения аргументов запроса \\\
	args := make([]interface{}, 0)
	argId := 1

	/// Проверки на наличие новых значений \\\
	if record.HospitalAddress != nil {
		values = append(values, fmt.Sprintf("hospital_address=$%d", argId))
		args = append(args, *record.HospitalAddress)
		argId++
	}
	if record.DoctorOffice != nil {
		values = append(values, fmt.Sprintf("doctor_office=$%d", argId))
		args = append(args, *record.DoctorOffice)
		argId++
	}
	if record.DoctorID != nil {
		values = append(values, fmt.Sprintf("doctor_id=$%d", argId))
		args = append(args, *record.DoctorID)
		argId++
	}
	if record.TimeRecord != nil {
		values = append(values, fmt.Sprintf("time_record=$%d", argId))
		args = append(args, *record.TimeRecord)
		argId++
	}

	/// Формирование строки со всеми измененными полями и их значениями \\\
	valuesQuery := strings.Join(values, ", ")
	query := fmt.Sprintf("UPDATE record  SET %s WHERE id = $%d", valuesQuery, argId)
	args = append(args, record.ID)

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), r.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	result, err := r.conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update record partially: %v", err)
	}

	if result.RowsAffected() == 0 {
		return apperror.ErrEmptyString
	}
	return nil
}

/// Функция DeleteRecord для сущности RecordStorage удаляет записи на прием из БД \\\

func (r *RecordStorage) DeleteRecord(id int64) error {
	r.logger.Info("POSTGRES: DELETE RECORD")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), r.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	result, err := r.conn.Exec(ctx,
		`DELETE FROM record WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete record: %v", err)
	}

	if result.RowsAffected() == 0 {
		return apperror.ErrEmptyString
	}
	return nil
}
