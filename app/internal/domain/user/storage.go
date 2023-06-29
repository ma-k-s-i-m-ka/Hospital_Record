package user

type Storage interface {
	Create(user *User) (*User, error)
	FindByEmail(email string) (*User, error)
	FindById(id int64) (*User, error)
	FindByPolicyNumber(policy string) (*User, error)
	Update(user *UpdateUserDTO) error
	PartiallyUpdate(user *PartiallyUpdateUserDTO) error
	Delete(id int64) error
}
