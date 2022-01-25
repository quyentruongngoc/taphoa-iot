package receipt

import (
	"taphoa-iot-backend/package/api-interface/receipt"
)

type Repository interface {
	AddReceipt(receipt.Instance, string) (receipt.Instance, error)
	UpdateReceipt(receipt.Instance, string) (receipt.Instance, error)
	DeleteReceipt(int64, string) error
	DescribeReceipt(int, string, string, int) (receipt.DescReceipt, error)
	ReportReceipt(string, string, string, int) (receipt.Report, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) AddReceipt(in receipt.Instance, user string) (receipt.Instance, error) {
	return s.repo.AddReceipt(in, user)
}

func (s *Service) ListReceipts(page int, user string, search string, status int) (receipt.DescReceipt, error) {
	return s.repo.DescribeReceipt(page, user, search, status)
}

func (s *Service) UpdateReceipt(in receipt.Instance, user string) (receipt.Instance, error) {
	return s.repo.UpdateReceipt(in, user)
}

func (s *Service) DeleteReceipt(id int64, user string) error {
	return s.repo.DeleteReceipt(id, user)
}

func (s *Service) ReportReceipt(user string, from string, to string, status int) (receipt.Report, error) {
	return s.repo.ReportReceipt(user, from, to, status)
}
