package receipt

import "taphoa-iot-backend/package/api-interface/warehouse"

type Instance struct {
	TotalMn int                  `json:"total_mn,omitempty"`
	Detail  []warehouse.Instance `json:"detail,omitempty"`
	Name    string               `json:"name,omitempty"`
	Phone   string               `json:"phone,omitempty"`
	Addr    string               `json:"addr,omitempty"`
	ID      uint64               `json:"id,omitempty"`
	Token   string               `json:"token,omitempty"`
	Status  int                  `json:"status,omitempty"`
}

type DescReceipt struct {
	TotalPage int        `json:"total_page,omitempty"`
	Page      int        `json:"page,omitempty"`
	Data      []Instance `json:"data,omitempty"`
}

type Report struct {
	Detail      []warehouse.Instance `json:"detail,omitempty"`
	TotalMn     int64                `json:"total_mn,omitempty"`
	TotalProfit int64                `json:"total_profit,omitempty"`
}

type API interface {
	AddReceipt(Instance, string) (Instance, error)
	UpdateReceipt(Instance, string) (Instance, error)
	DeleteReceipt(int64, string) error
	ListReceipts(int, string, string, int) (DescReceipt, error)
	ReportReceipt(string, string, string, int) (Report, error)
}
