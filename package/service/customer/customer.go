package customer

import "taphoa-iot-backend/package/api-interface/customer"

type Repository interface {
	DescribeCustom(int, string, string) (customer.DescCustom, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) ListCustomers(page int, user string, search string) (customer.DescCustom, error) {
	return s.repo.DescribeCustom(page, user, search)
}
