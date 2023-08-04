package user

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

var _ Storage = &UserStorage{}

/// Структура UserStorage содержащая поля для работы с БД \\\

type UserStorage struct {
	logger         logger.Logger
	conn           *pgx.Conn
	requestTimeout time.Duration
}

/// Структура NewStorage возвращает новый экземпляр UserStorage инициализируя переданные в него аргументы \\\

func NewStorage(storage *pgx.Conn, requestTimeout int) Storage {
	return &UserStorage{
		logger:         logger.GetLogger(),
		conn:           storage,
		requestTimeout: time.Duration(requestTimeout) * time.Second,
	}
}

/// Функция Create для сущности UserStorage создает записи пациентов в БД \\\

func (d *UserStorage) Create(user *User) (*User, error) {
	d.logger.Info("POSTGRES: CREATE USER")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	row := d.conn.QueryRow(ctx,
		`INSERT INTO patients (email, name, surname, age, gender, password, policy_number)
			 VALUES($1,$2,$3,$4,$5,$6,$7) 
			 RETURNING id, created_at`,
		user.Email, user.Name, user.Surname, user.Age, user.Gender, user.Password, user.PolicyNumber)

	/// Сканирование полученных значений из БД \\\
	err := row.Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		err = fmt.Errorf("failed to execute create user query: %v", err)
		d.logger.Error(err)
		return nil, err
	}
	return user, nil
}

/// Функция FindByEmail для сущности UserStorage получает записи пациентов из БД по адресу электронной почты \\\

func (d *UserStorage) FindByEmail(email string) (*User, error) {
	d.logger.Info("POSTGRES: GET USER BY EMAIL")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	row := d.conn.QueryRow(ctx,
		`SELECT * FROM patients
			 WHERE email = $1`, email)
	user := &User{}

	/// Сканирование полученных значений из БД \\\
	err := row.Scan(
		&user.ID, &user.Email, &user.Name, &user.Surname, &user.Patronymic,
		&user.Age, &user.Gender, &user.PhoneNumber, &user.Address, &user.Password,
		&user.PolicyNumber, &user.DiseaseID, &user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		err = fmt.Errorf("failed to execute find user by email query: %v", err)
		d.logger.Error(err)
		return nil, err
	}
	return user, nil
}

/// Функция FindById для сущности UserStorage получает записи пациента из БД по id \\\

func (d *UserStorage) FindById(id int64) (*User, error) {
	d.logger.Info("POSTGRES: GET USER BY ID")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	row := d.conn.QueryRow(ctx,
		`SELECT * FROM patients
			 WHERE id = $1`, id)
	user := &User{}

	/// Сканирование полученных значений из БД \\\
	err := row.Scan(
		&user.ID, &user.Email, &user.Name, &user.Surname, &user.Patronymic,
		&user.Age, &user.Gender, &user.PhoneNumber, &user.Address, &user.Password,
		&user.PolicyNumber, &user.DiseaseID, &user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrEmptyString
		}
		err = fmt.Errorf("failed to execute find user by id query: %v", err)
		d.logger.Error(err)
		return nil, err
	}
	return user, nil
}

/// Функция FindByPolicyNumber для сущности UserStorage получает записи пациентов из БД по номеру полиса \\\

func (d *UserStorage) FindByPolicyNumber(policy string) (*User, error) {
	d.logger.Info("POSTGRES: GET USER BY POLICY NUMBER")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	row := d.conn.QueryRow(ctx,
		`SELECT * FROM patients
			 WHERE policy_number = $1`, policy)
	user := &User{}

	/// Сканирование полученных значений из БД \\\
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.Surname, &user.Patronymic,
		&user.Age, &user.Gender, &user.PhoneNumber, &user.Address, &user.Password,
		&user.PolicyNumber, &user.DiseaseID, &user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		err = fmt.Errorf("failed to execute find user by policynumber query: %v", err)
		d.logger.Error(err)
		return nil, err
	}
	return user, nil
}

/// Функция Update для сущности UserStorage обновляет записи о пациенте в БД \\\

func (d *UserStorage) Update(user *UpdateUserDTO) error {
	d.logger.Info("POSTGRES: UPDATE USER")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	result, err := d.conn.Exec(ctx,
		`UPDATE patients
			 SET email=$1, name=$2, surname=$3, patronymic=$4, age=$5, gender=$6, phone_number=$7, address=$8, password=$9, policy_number=$10, disease_id=$11
			 WHERE id =$12`,
		user.Email, user.Name, user.Surname, user.Patronymic, user.Age, user.Gender, user.PhoneNumber, user.Address, user.Password, user.PolicyNumber, user.DiseaseID, user.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.ErrEmptyString
		}
		err = fmt.Errorf("failed to execute update user query: %v", err)
		d.logger.Error(err)
		return err
	}
	if result.RowsAffected() == 0 {
		return apperror.ErrEmptyString
	}
	return nil
}

/// Функция PartiallyUpdate для сущности UserStorage частично обновляет записи о пациенте в БД \\\

func (d *UserStorage) PartiallyUpdate(user *PartiallyUpdateUserDTO) error {
	d.logger.Info("POSTGRES: PARTIALLY UPDATE USER")

	/// Создание пустого слайса для хранения обновляемых строк \\\
	values := make([]string, 0)

	/// Создание пустого слайса для хранения аргументов запроса \\\
	args := make([]interface{}, 0)
	argId := 1

	/// Проверки на наличие новых значений \\\
	if user.Email != nil {
		values = append(values, fmt.Sprintf("email=$%d", argId))
		args = append(args, *user.Email)
		argId++
	}
	if user.PhoneNumber != nil {
		values = append(values, fmt.Sprintf("phone_number=$%d", argId))
		args = append(args, *user.PhoneNumber)
		argId++
	}
	if user.Address != nil {
		values = append(values, fmt.Sprintf("address=$%d", argId))
		args = append(args, *user.Address)
		argId++
	}
	if user.Password != nil {
		values = append(values, fmt.Sprintf("password=$%d", argId))
		args = append(args, *user.Password)
		argId++
	}
	if user.DiseaseID != nil {
		values = append(values, fmt.Sprintf("disease_id=$%d", argId))
		args = append(args, *user.DiseaseID)
		argId++
	}

	/// Формирование строки со всеми измененными полями и их значениями \\\
	valuesQuery := strings.Join(values, ", ")
	query := fmt.Sprintf("UPDATE patients  SET %s WHERE id = $%d", valuesQuery, argId)
	args = append(args, user.ID)

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	result, err := d.conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user partially: %v", err)
	}

	if result.RowsAffected() == 0 {
		return apperror.ErrEmptyString
	}
	return nil
}

/// Функция Delete для сущности UserStorage удаляет записи о пациентах из БД \\\

func (d *UserStorage) Delete(id int64) error {
	d.logger.Info("POSTGRES: DELETE USER")

	/// Ограничение времени выполнения запроса \\\
	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	/// Выполнение запроса к БД \\\
	result, err := d.conn.Exec(ctx,
		`DELETE FROM patients WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	if result.RowsAffected() == 0 {
		return apperror.ErrEmptyString
	}
	return nil
}
