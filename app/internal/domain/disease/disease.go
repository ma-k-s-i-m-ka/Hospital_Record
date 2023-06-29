package disease

type Disease struct {
	ID          int64  `json:"id" example:"1567"`
	BodyPart    string `json:"body_part" example:"hand"`
	Description string `json:"description" example:"broken finger"`
}

type CreateDiseaseDTO struct {
	BodyPart    string `json:"body_part" example:"hand"`
	Description string `json:"description" example:"broken finger"`
}
type UpdateDiseaseDTO struct {
	ID          int64  `json:"id" example:"1567"`
	BodyPart    string `json:"body_part" example:"hand"`
	Description string `json:"description" example:"broken finger"`
}
