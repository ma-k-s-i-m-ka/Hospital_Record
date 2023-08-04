package auth

import "golang.org/x/crypto/bcrypt"

/// Структура для авторизации и регистрации пользователей \\\

type AccessToken struct {
	ID           int64   `json:"id" example:"1567"`
	Email        string  `json:"email" example:"petrovmaksim1992@mail.ru"`
	Name         string  `json:"name" example:"Maksim"`
	Surname      string  `json:"surname" example:"Petrov"`
	Patronymic   *string `json:"patronymic,omitempty" example:"Olegovich"`
	Age          uint8   `json:"age" example:"20"`
	Gender       string  `json:"gender" example:"female"`
	PhoneNumber  *string `json:"phone_number,omitempty" example:"85555555555"`
	Address      *string `json:"address,omitempty" example:"Moscow, Yaroslavskoe shosse, 26 korpus 12"`
	PolicyNumber string  `json:"policy_number" example:"2197799730000060"`
}

type RefreshToken struct {
	ID int64 `json:"id" example:"1567"`
}

type AuthByEmail struct {
	Email    string `json:"email" example:"petrovmaksim1992@mail.ru"`
	Password string `json:"password" example:"abcdEFG"`
}

type AuthByPolicyNumber struct {
	PolicyNumber string `json:"policy_number" example:"2197799730000060"`
	Password     string `json:"password" example:"abcdEFG"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Register struct {
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
type RegisterResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (u *Register) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
