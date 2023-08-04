package user

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

/// Структура для создания и обновления пациентов \\\

type User struct {
	ID           int64     `json:"id" example:"1567"`
	Email        string    `json:"email" example:"petrovmaksim1992@mail.ru"`
	Name         string    `json:"name" example:"Maksim"`
	Surname      string    `json:"surname" example:"Petrov"`
	Patronymic   *string   `json:"patronymic,omitempty" example:"Olegovich"`
	Age          uint8     `json:"age" example:"20"`
	Gender       string    `json:"gender" example:"female"`
	PhoneNumber  *string   `json:"phone_number,omitempty" example:"85555555555"`
	Address      *string   `json:"address,omitempty" example:"Moscow, Yaroslavskoe shosse, 26 korpus 12"`
	Password     string    `json:"password"`
	PolicyNumber string    `json:"policy_number" example:"2197799730000060"`
	DiseaseID    *[]int64  `json:"disease_id,omitempty" example:"123"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

type CreateUserDTO struct {
	Email        string  `json:"email" example:"petrovmaksim1992@mail.ru"`
	Name         string  `json:"name" example:"Maksim"`
	Surname      string  `json:"surname" example:"Petrov"`
	Patronymic   *string `json:"patronymic,omitempty" example:"Olegovich"`
	Age          uint8   `json:"age" example:"20"`
	Gender       string  `json:"gender" example:"female"`
	PhoneNumber  *string `json:"phone_number,omitempty" example:"85555555555"`
	Address      *string `json:"address,omitempty" example:"Moscow, Yaroslavskoe shosse, 26 korpus 12"`
	Password     string  `json:"password" example:"sfdsg"`
	PolicyNumber string  `json:"policy_number" example:"2197799730000060"`
}

type UpdateUserDTO struct {
	ID           int64   `json:"id"`
	Email        string  `json:"email" example:"petrovmaksim1992@mail.ru"`
	Name         string  `json:"name" example:"Maksim"`
	Surname      string  `json:"surname" example:"Petrov"`
	Patronymic   *string `json:"patronymic,omitempty" example:"Olegovich"`
	Age          uint8   `json:"age" example:"20"`
	Gender       string  `json:"gender" example:"female"`
	PhoneNumber  *string `json:"phone_number,omitempty" example:"85555555555"`
	Address      *string `json:"address,omitempty" example:"Moscow, Yaroslavskoe shosse, 26 korpus 12"`
	Password     string  `json:"password" example:"abcdEFG"`
	PolicyNumber string  `json:"policy_number" example:"2197799730000060"`
	DiseaseID    []int64 `json:"disease_id,omitempty" example:"123"`
}
type PartiallyUpdateUserDTO struct {
	ID          int64    `json:"id"`
	Email       *string  `json:"email" example:"petrovmaksim1992@mail.ru"`
	PhoneNumber *string  `json:"phone_number,omitempty" example:"85555555555"`
	Address     *string  `json:"address,omitempty" example:"Moscow, Yaroslavskoe shosse, 26 korpus 12"`
	Password    *string  `json:"password" example:"abcdEFG"`
	DiseaseID   *[]int64 `json:"disease_id,omitempty" example:"123"`
}

/// Хэширование паролей \\\

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *UpdateUserDTO) HashPassword() error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *PartiallyUpdateUserDTO) HashPassword() error {
	if u.Password == nil {
		return fmt.Errorf("new password cannot be nil")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	*u.Password = string(hashedPassword)
	return nil
}

/// Проверка введенного пароля на соответсвие паролю пользователя в БД \\\

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
