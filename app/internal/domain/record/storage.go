package record

type Storage interface {
	CreateRecord(record *Record) (*Record, error)
	FindRecordByPatientsId(id int64) (*Record, error)
	FindRecordById(id int64) (*Record, error)
	UpdateRecord(record *UpdateRecordDTO) error
	PartiallyUpdateRecord(record *PartiallyUpdateRecordDTO) error
	DeleteRecord(id int64) error
}
