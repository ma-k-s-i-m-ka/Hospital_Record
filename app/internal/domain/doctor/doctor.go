package doctor

type Doctor struct {
	ID                   int64   `json:"id" example:"1567"`
	Name                 string  `json:"name" example:"Vitaliy"`
	Surname              string  `json:"surname" example:"Ivanov"`
	Patronymic           *string `json:"patronymic,omitempty" example:"Semenovich"`
	ImageID              int64   `json:"image_id" example:"1567"`
	Gender               string  `json:"gender" example:"male"`
	Rating               float32 `json:"rating" example:"4"`
	Age                  int8    `json:"age" example:"28"`
	RecordingIsAvailable bool    `json:"recording_is_available" example:"true"`
	SpecializationID     int64   `json:"specialization_id" example:"1567"`
	PortfolioID          int64   `json:"portfolio_id" example:"1567"`
}

type CreateDoctorDTO struct {
	Name                 string  `json:"name" example:"Vitaliy"`
	Surname              string  `json:"surname" example:"Ivanov"`
	Patronymic           *string `json:"patronymic,omitempty" example:"Semenovich"`
	ImageID              int64   `json:"image_id" example:"1567"`
	Gender               string  `json:"gender" example:"male"`
	Rating               float32 `json:"rating" example:"4"`
	Age                  int8    `json:"age" example:"28"`
	RecordingIsAvailable bool    `json:"recording_is_available" example:"true"`
	SpecializationID     int64   `json:"specialization_id" example:"1567"`
	PortfolioID          int64   `json:"portfolio_id" example:"1567"`
}

type UpdateDoctorDTO struct {
	ID                   int64   `json:"id"`
	Name                 string  `json:"name" example:"Vitaliy"`
	Surname              string  `json:"surname" example:"Ivanov"`
	Patronymic           *string `json:"patronymic" example:"Semenovich"`
	ImageID              int64   `json:"image_id" example:"1567"`
	Gender               string  `json:"gender" example:"male"`
	Rating               float32 `json:"rating" example:"4"`
	Age                  uint8   `json:"age" example:"28"`
	RecordingIsAvailable bool    `json:"recording_is_available" example:"true"`
	SpecializationID     int64   `json:"specialization_id" example:"1567"`
	PortfolioID          int64   `json:"portfolio_id" example:"1567"`
}

type PartiallyUpdateDoctorDTO struct {
	ID                   int64    `json:"id"`
	ImageID              *int64   `json:"image_id" example:"1567"`
	Rating               *float32 `json:"rating" example:"4"`
	RecordingIsAvailable *bool    `json:"recording_is_available" example:"true"`
}
