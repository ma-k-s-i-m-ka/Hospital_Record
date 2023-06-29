package user

type User struct {
	ID           int64   `json:"id" example:"1567"`
	Email        string  `json:"email" example:"petrovmaksim1992@mail.ru"`
	Name         string  `json:"name" example:"Maksim"`
	Surname      string  `json:"surname" example:"Petrov"`
	Patronymic   *string `json:"patronymic,omitempty" example:"Olegovich"`
	Age          uint8   `json:"age" example:"20"`
	Gender       string  `json:"gender" example:"female"`
	PhoneNumber  *string `json:"phone_number,omitempty" example:"85555555555"`
	Address      *string `json:"address,omitempty" example:"Moscow, Yaroslavskoe shosse, 26 korpus 12"`
	Password     string  `json:"-"`
	PolicyNumber string  `json:"policy_number" example:"2197799730000060"`
	DiseaseID    *int64  `json:"disease_id,omitempty" example:"123"`
	CreatedAt    string  `json:"created_at,omitempty"`
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
	Password     string  `json:"-" example:"-"`
	PolicyNumber string  `json:"policy_number" example:"2197799730000060"`
}

type UpdateUserDTO struct {
	ID           int64   `json:"-"`
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
	DiseaseID    *int64  `json:"disease_id,omitempty" example:"123"`
}
type PartiallyUpdateUserDTO struct {
	ID          int64   `json:"-"`
	Email       *string `json:"email" example:"petrovmaksim1992@mail.ru"`
	PhoneNumber *string `json:"phone_number,omitempty" example:"85555555555"`
	Address     *string `json:"address,omitempty" example:"Moscow, Yaroslavskoe shosse, 26 korpus 12"`
	Password    *string `json:"password" example:"abcdEFG"`
	DiseaseID   *int64  `json:"disease_id,omitempty" example:"123"`
}
