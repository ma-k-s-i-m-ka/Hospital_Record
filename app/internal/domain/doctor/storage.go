package doctor

type Storage interface {
	Create(doctor *Doctor) (*Doctor, error)
	FindAll() ([]Doctor, error)
	FindAllAvailable(id int64, recordingIsAvailable bool) ([]Doctor, error)
	FindById(id int64) (*Doctor, error)
	FindByPortfolioId(id int64) (*Doctor, error)
	Update(doctor *UpdateDoctorDTO) error
	PartiallyUpdate(doctor *PartiallyUpdateDoctorDTO) error
	Delete(id int64) error
}
