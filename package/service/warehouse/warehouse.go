package warehouse

import "taphoa-iot-backend/package/api-interface/warehouse"

type Repository interface {
	CreateItem(warehouse.Instance, string) (warehouse.Instance, error)
	UpdateItem(warehouse.Instance, string) (warehouse.Instance, error)
	DescribeItems(int, string, string) (warehouse.DescItems, error)
	GetItem(uint64, string) (warehouse.Instance, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) CreateItem(in warehouse.Instance, user string) (warehouse.Instance, error) {
	return s.repo.CreateItem(in, user)
}

func (s *Service) DescribeItems(page int, user string, search string) (warehouse.DescItems, error) {
	return s.repo.DescribeItems(page, user, search)
}

func (s *Service) UpdateItem(in warehouse.Instance, user string) (warehouse.Instance, error) {
	data, err := s.repo.GetItem(in.ID, user)
	if err != nil {
		return warehouse.Instance{}, err
	}
	in.Remain = data.Remain + in.Total

	return s.repo.UpdateItem(in, user)
}
