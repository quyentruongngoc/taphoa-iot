package repository

import (
	"fmt"
	"log"
	"taphoa-iot-backend/package/api-interface/warehouse"
)

const (
	ITEM_PER_PAGE = 20
)

type Warehouse struct {
	BaseModel

	User          string `sql:"type:varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci"`
	Name          string `sql:"type:text CHARACTER SET utf8 COLLATE utf8_general_ci"`
	PurchasePrice int
	SalePrice     int
	Total         int64
	Profit        int64
	Sold          int64
	Remain        int64
	SaleTotal     int
}

func (s *Storage) CreateItem(in warehouse.Instance, user string) (warehouse.Instance, error) {
	record := Warehouse{
		User:          user,
		Name:          in.Name,
		PurchasePrice: in.PurchasePrice,
		SalePrice:     in.SalePrice,
		Total:         in.Total,
		Remain:        in.Total,
		Sold:          0,
	}

	result := s.db.Model(&record).Debug().Create(&record)
	if result.Error != nil {
		return warehouse.Instance{}, result.Error
	}

	return in, nil
}

func (s *Storage) UpdateItem(in warehouse.Instance, user string) (warehouse.Instance, error) {
	if in.ID == 0 {
		return warehouse.Instance{}, fmt.Errorf("ID should not zero")
	}

	record := Warehouse{
		Name:          in.Name,
		PurchasePrice: in.PurchasePrice,
		SalePrice:     in.SalePrice,
		Remain:        in.Remain,
	}

	result := db.Model(&record).Debug().Where(&Warehouse{
		BaseModel: BaseModel{
			ID: in.ID,
		},
		User: user,
	}).Updates(&record)
	if result.Error != nil {
		return warehouse.Instance{}, result.Error
	}

	return in, nil
}

func (s *Storage) GetItem(ID uint64, user string) (warehouse.Instance, error) {
	var record Warehouse

	where := Warehouse{
		BaseModel: BaseModel{
			ID: ID,
		},
		User: user,
	}

	result := s.db.Debug().Model(&where).Where(&where).First(&record)
	if result.Error != nil {
		return warehouse.Instance{}, result.Error
	}

	return warehouse.Instance{
		ID:            ID,
		Name:          record.Name,
		PurchasePrice: record.PurchasePrice,
		SalePrice:     record.SalePrice,
		Total:         record.Total,
		Sold:          record.Sold,
		Remain:        record.Remain,
		Token:         "",
	}, nil

}

func (s *Storage) DescribeItems(page int, user string, search string) (warehouse.DescItems, error) {
	var records []Warehouse
	var ret warehouse.DescItems
	var dataArray []warehouse.Instance

	search = "%" + search + "%"

	result := s.db.Debug().Model(&Warehouse{}).
		Where("user = ? AND name LIKE ?", user, search).Find(&records)
	if result.Error != nil {
		return warehouse.DescItems{}, nil
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
		Order("remain desc, id").
		Offset(offset).
		Limit(ITEM_PER_PAGE).
		Find(&records)
	if result.Error != nil {
		return warehouse.DescItems{}, nil
	}

	for _, rec := range records {
		ins := warehouse.Instance{
			ID:            rec.ID,
			Name:          rec.Name,
			PurchasePrice: rec.PurchasePrice,
			SalePrice:     rec.SalePrice,
			Total:         rec.Total,
			Sold:          rec.Sold,
			Remain:        rec.Remain,
			Token:         "",
		}
		dataArray = append(dataArray, ins)
	}
	ret.Data = dataArray

	return ret, nil
}
