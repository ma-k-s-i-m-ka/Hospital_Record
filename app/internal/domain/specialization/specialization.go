package specialization

type Specialization struct {
	ID   int64  `json:"id" example:"1567"`
	Name string `json:"name_specialization" example:"therapist"`
}

type CreateSpecializationDTO struct {
	Name string `json:"name_specialization" example:"therapist"`
}

type UpdateSpecializationDTO struct {
	ID   int64  `json:"id" example:"1567"`
	Name string `json:"name_specialization" example:"therapist"`
}
