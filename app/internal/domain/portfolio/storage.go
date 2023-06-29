package portfolio

type Storage interface {
	Create(portfolio *Portfolio) (*Portfolio, error)
	FindById(id int64) (*Portfolio, error)
	Update(portfolio *UpdatePortfolioDTO) error
	Delete(id int64) error
}
