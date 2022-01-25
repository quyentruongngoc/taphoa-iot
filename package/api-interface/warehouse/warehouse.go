package warehouse

type Instance struct {
	ID            uint64 `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	PurchasePrice int    `json:"purchase_price,omitempty"`
	SalePrice     int    `json:"sale_price,omitempty"`
	Total         int64  `json:"total,omitempty"`
	Sold          int64  `json:"sold,omitempty"`
	Remain        int64  `json:"remain,omitempty"`
	Token         string `json:"token,omitempty"`
	Quantity      int    `json:"quantity,omitempty"`
	Profit        int64  `json:"profit,omitempty"`
	SaleTotal     int    `json:"sale_total,omitempty"`
}

type DescItems struct {
	TotalPage int        `json:"total_page,omitempty"`
	Page      int        `json:"page,omitempty"`
	Data      []Instance `json:"data,omitempty"`
}

type API interface {
	CreateItem(Instance, string) (Instance, error)
	UpdateItem(Instance, string) (Instance, error)
	DescribeItems(int, string, string) (DescItems, error)
}
