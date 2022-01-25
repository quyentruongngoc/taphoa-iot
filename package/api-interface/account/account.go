package account

type Instance struct {
	User   string `json:"user,omitempty"`
	Passwd string `json:"passwd,omitempty"`
	Role   int    `json:"role,omitempty"`
	Token  string `json:"token,omitempty"`
	Name   string `json:"name,omitempty"`
	Addr   string `json:"addr,omitempty"`
	Phone  string `json:"phone,omitempty"`
	Email  string `json:"email,omitempty"`
	Gender bool   `json:"gender,omitempty"`
}

type API interface {
	Create(Instance) (Instance, error)
	Authenticate(Instance) (Instance, error)
}
