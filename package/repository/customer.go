package repository

import (
	"log"
	"taphoa-iot-backend/package/api-interface/customer"
)

type Customer struct {
	BaseModel
	Phone string `sql:"type:varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci;unique;not null"`
	User  string `sql:"type:varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci"`
	Name  string `sql:"type:varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci"`
	Addr  string `sql:"type:varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci"`
	Total int64
}

func (s *Storage) DescribeCustom(page int, user string, search string) (customer.DescCustom, error) {
	var records []Customer
	var ret customer.DescCustom
	var dataArray []customer.Instance

	search = "%" + search + "%"

	result := s.db.Debug().Model(&Customer{}).
		Where("user = ? AND phone LIKE ?", user, search).Find(&records)
	if result.Error != nil {
		return customer.DescCustom{}, nil
	}

	totalRec := int(result.RowsAffected)
	ret.TotalPage = totalRec / ITEM_PER_PAGE
	log.Printf("Quyen debug mod: %v", (totalRec % ITEM_PER_PAGE))
	log.Printf("Quyen debug RowsAffected: %v", totalRec)
	if (int(result.RowsAffected) % ITEM_PER_PAGE) != 0 {
		ret.TotalPage += 1
	}
	ret.Page = page + 1
	offset := page * ITEM_PER_PAGE

	result = result.
		Order("total desc, id").
		Offset(offset).
		Limit(ITEM_PER_PAGE).
		Find(&records)
	if result.Error != nil {
		return customer.DescCustom{}, nil
	}

	for _, rec := range records {
		ins := customer.Instance{
			User:  rec.User,
			Token: "",
			Name:  rec.Name,
			Addr:  rec.Addr,
			Phone: rec.Phone,
			Total: rec.Total,
		}
		dataArray = append(dataArray, ins)
	}
	ret.Data = dataArray

	return ret, nil
}

func (s *Storage) GetCustomer(user string, phone string) (Customer, error) {
	var record Customer

	result := s.db.Model(&record).Debug().Where(&Customer{
		Phone: phone,
	}).First(&record)
	if result.Error != nil {
		return Customer{}, result.Error
	}

	return record, nil
}

func (s *Storage) CreateOrUpdateCustomer(user string, custom customer.Instance) error {
	record := Customer{
		Phone: custom.Phone,
		User:  user,
		Name:  custom.Name,
		Addr:  custom.Addr,
	}

	temp, _ := s.GetCustomer(user, custom.Phone)
	// if err != nil {
	// 	log.Printf("Quyen debug: failed to get customer: %v\n", err)
	// 	return err
	// }

	record.Total = temp.Total + custom.Total

	if db.Model(&record).Debug().Where(&Customer{
		User:  user,
		Phone: custom.Phone,
	}).Updates(&record).RowsAffected == 0 {
		db.Debug().Create(&record)
	}

	return nil
}
