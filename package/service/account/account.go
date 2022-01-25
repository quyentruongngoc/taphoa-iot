package account

import (
	"fmt"
	"log"
	"taphoa-iot-backend/internal"
	"taphoa-iot-backend/package/api-interface/account"
	"time"
)

type Repository interface {
	CreateUser(account.Instance) error
	DescribeUser(string) (account.Instance, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) Create(in account.Instance) (account.Instance, error) {
	if len(in.Passwd) == 0 {
		return account.Instance{}, fmt.Errorf("Password should not empty")
	}

	if len(in.Passwd) > 0 {
		in.Passwd = internal.SHA256(in.Passwd)
	} else {
		in.Passwd = ""
	}

	err := s.repo.CreateUser(in)
	if err != nil {
		return account.Instance{}, err
	}

	return in, nil
}

func (s *Service) Authenticate(in account.Instance) (account.Instance, error) {
	instance, err := s.repo.DescribeUser(in.User)
	if err != nil {
		log.Printf("Failed to get data from database: %v", err)
		return account.Instance{}, err
	}
	hashPw := internal.SHA256(in.Passwd)

	if instance.Passwd == hashPw {
		now := time.Now().Nanosecond()
		temp := fmt.Sprintf("%v-%v", instance.User, now)
		log.Printf("Quyen debug Token before hash: %v\n", temp)
		instance.Token = internal.SHA256(temp)
		log.Printf("Generate new Token: %v\n", instance.Token)
		// update token to list
		internal.AddToTokenList(instance.Token, uint(instance.Role), instance.User)

		return account.Instance{
			User:   instance.User,
			Role:   instance.Role,
			Passwd: "",
			Token:  instance.Token,
		}, nil
	}

	return account.Instance{}, fmt.Errorf("Password not correct")
}
