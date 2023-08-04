package record

import "time"

/// Структура для создания и обновления записей \\\

type Record struct {
	ID               int64     `json:"id" example:"1567"`
	HospitalAddress  string    `json:"hospital_address" example:"Roterta, dom 12"`
	DoctorOffice     string    `json:"doctor_office" example:"201B"`
	Tagging          string    `json:"tagging" example:"Nichego ne est za 3 chasa pered priemom"`
	PatientsID       int64     `json:"patients_id" example:"1"`
	DoctorID         int64     `json:"doctor_id" example:"1"`
	SpecializationID int64     `json:"specialization_id" example:"1"`
	TimeRecord       time.Time `json:"time_record" example:"2023-07-27T15:30:00Z"`
}

type CreateRecordDTO struct {
	ID               int64     `json:"id" example:"1567"`
	HospitalAddress  string    `json:"hospital_address" example:"Roterta, dom 12"`
	DoctorOffice     string    `json:"doctor_office" example:"201B"`
	Tagging          string    `json:"tagging" example:"Nichego ne est za 3 chasa pered priemom"`
	PatientsID       int64     `json:"patients_id" example:"1"`
	DoctorID         int64     `json:"doctor_id" example:"1"`
	SpecializationID int64     `json:"specialization_id" example:"1"`
	TimeRecord       time.Time `json:"time_record" example:"2023-07-27T15:30:00Z"`
}
type UpdateRecordDTO struct {
	ID               int64     `json:"id" example:"1567"`
	HospitalAddress  string    `json:"hospital_address" example:"Roterta, dom 12"`
	DoctorOffice     string    `json:"doctor_office" example:"201B"`
	Tagging          string    `json:"tagging" example:"Nichego ne est za 3 chasa pered priemom"`
	PatientsID       int64     `json:"patients_id" example:"1"`
	DoctorID         int64     `json:"doctor_id" example:"1"`
	SpecializationID int64     `json:"specialization_id" example:"1"`
	TimeRecord       time.Time `json:"time_record" example:"2023-07-27T15:30:00Z"`
}

type PartiallyUpdateRecordDTO struct {
	ID              int64      `json:"id" example:"1567"`
	HospitalAddress *string    `json:"hospital_address" example:"Roterta, dom 12"`
	DoctorOffice    *string    `json:"doctor_office" example:"201B"`
	DoctorID        *int64     `json:"doctor_id" example:"1"`
	TimeRecord      *time.Time `json:"time_record" example:"2023-07-27T15:30:00Z"`
}
