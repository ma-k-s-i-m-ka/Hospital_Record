package portfolio

type Portfolio struct {
	ID             int64  `json:"id" example:"1567"`
	Education      string `json:"education" example:"residency Institute of N. I. Pirogov"`
	Awards         string `json:"awards,omitempty" example:"Advanced training course of 4 categories"`
	WorkExperience uint8  `json:"work_experience" example:"25"`
}

type CreatePortfolioDTO struct {
	Education      string `json:"education" example:"residency Institute of N. I. Pirogov"`
	Awards         string `json:"awards,omitempty" example:"Advanced training course of 4 categories"`
	WorkExperience uint8  `json:"work_experience" example:"25"`
}

type UpdatePortfolioDTO struct {
	ID             int64  `json:"id" example:"1567"`
	Education      string `json:"education" example:"residency Institute of N. I. Pirogov"`
	Awards         string `json:"awards,omitempty" example:"Advanced training course of 4 categories"`
	WorkExperience uint8  `json:"work_experience" example:"25"`
}
