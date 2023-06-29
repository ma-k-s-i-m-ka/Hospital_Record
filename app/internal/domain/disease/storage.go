package disease

type Storage interface {
	Create(disease *Disease) (*Disease, error)
	FindById(id int64) (*Disease, error)
	Update(disease *UpdateDiseaseDTO) error
	Delete(id int64) error
}
