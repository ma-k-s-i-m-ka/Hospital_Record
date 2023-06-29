package specialization

type Storage interface {
	Create(specialization *Specialization) (*Specialization, error)
	FindById(id int64) (*Specialization, error)
	Update(specialization *UpdateSpecializationDTO) error
	Delete(id int64) error
}
