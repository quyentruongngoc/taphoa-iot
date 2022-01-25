package customer

type Instance struct {
	User  string `json:"user,omitempty"`
	Token string `json:"token,omitempty"`
	Name  string `json:"name,omitempty"`
	Addr  string `json:"addr,omitempty"`
	Phone string `json:"phone,omitempty"`
	Total int64  `json:"total,omitempty"`
}

type DescCustom struct {
	TotalPage int        `json:"total_page,omitempty"`
	Page      int        `json:"page,omitempty"`
	Data      []Instance `json:"data,omitempty"`
}

type API interface {
	ListCustomers(int, string, string) (DescCustom, error)
}
