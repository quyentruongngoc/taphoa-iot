package repository

import (
	"encoding/json"
	"log"
	"taphoa-iot-backend/package/api-interface/customer"
	"taphoa-iot-backend/package/api-interface/receipt"
	"taphoa-iot-backend/package/api-interface/warehouse"
	"time"

	"github.com/jinzhu/gorm"
)

type Receipt struct {
	BaseModel

	TotalMn int
	Detail  string `sql:"type:text CHARACTER SET utf8 COLLATE utf8_general_ci"`
	Name    string `sql:"type:varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci"`
	Phone   string `sql:"type:varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci"`
	Adress  string `sql:"type:text CHARACTER SET utf8 COLLATE utf8_general_ci"`
	User    string `sql:"type:varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci"`
	Status  int
}

func (s *Storage) AddReceipt(in receipt.Instance, user string) (receipt.Instance, error) {
	detailString, err := json.Marshal(in.Detail)
	if err != nil {
		return receipt.Instance{}, err
	}

	// check all item exit
	for _, item := range in.Detail {
		_, err := s.GetItem(item.ID, user)
		if err != nil {
			return receipt.Instance{}, err
		}
	}

	// update remain
	var totalmn int
	totalmn = 0
	for _, item := range in.Detail {
		totalmn += item.Quantity * item.SalePrice
		ins, _ := s.GetItem(item.ID, user)
		ins.Remain = ins.Remain - int64(item.Quantity)
		s.UpdateItem(ins, user)
	}

	record := Receipt{
		TotalMn: totalmn,
		Detail:  string(detailString),
		Name:    in.Name,
		Phone:   in.Phone,
		Adress:  in.Addr,
		User:    user,
		Status:  in.Status,
	}

	result := s.db.Model(&record).Debug().Create(&record)
	if result.Error != nil {
		return receipt.Instance{}, result.Error
	}

	custom := customer.Instance{
		User:  user,
		Token: "",
		Name:  in.Name,
		Addr:  in.Addr,
		Phone: in.Phone,
		Total: int64(totalmn),
	}
	s.CreateOrUpdateCustomer(user, custom)

	return in, nil
}

func (s *Storage) GetReceipt(id uint64, user string) (receipt.Instance, error) {
	var record Receipt
	var ret receipt.Instance

	result := s.db.Debug().Model(&Receipt{}).Where(
		&Receipt{
			BaseModel: BaseModel{
				ID: id,
			},
			User: user,
		}).First(&record)
	if result.Error != nil {
		return receipt.Instance{}, nil
	}

	ret = receipt.Instance{
		ID:      record.ID,
		TotalMn: record.TotalMn,
		Name:    record.Name,
		Phone:   record.Phone,
		Addr:    record.Adress,
		Token:   "",
		Status:  record.Status,
	}
	log.Printf("Quyen debug detail data json: %+v", record.Detail)
	detailData := []warehouse.Instance{}
	json.Unmarshal([]byte(record.Detail), &detailData)
	log.Printf("Quyen debug detail data %+v", detailData)
	ret.Detail = detailData

	return ret, nil
}

func (s *Storage) DescribeReceipt(page int, user string, search string, status int) (receipt.DescReceipt, error) {
	var records []Receipt
	var ret receipt.DescReceipt
	var dataArray []receipt.Instance

	search = "%" + search + "%"

	var result *gorm.DB

	if status == 255 {
		result = s.db.Debug().Model(&Receipt{}).
			Where("user = ? AND name LIKE ?", user, search).Find(&records)
		if result.Error != nil {
			return receipt.DescReceipt{}, nil
		}
	} else {
		result = s.db.Debug().Model(&Receipt{}).
			Where("user = ? AND name LIKE ? AND status = ?", user, search, status).Find(&records)
		if result.Error != nil {
			return receipt.DescReceipt{}, nil
		}
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
		Order("id").
		Offset(offset).
		Limit(ITEM_PER_PAGE).
		Find(&records)
	if result.Error != nil {
		return receipt.DescReceipt{}, nil
	}

	for _, rec := range records {
		ins := receipt.Instance{
			ID:      rec.ID,
			TotalMn: rec.TotalMn,
			Name:    rec.Name,
			Phone:   rec.Phone,
			Addr:    rec.Adress,
			Token:   "",
			Status:  rec.Status,
		}
		log.Printf("Quyen debug detail data json: %+v", rec.Detail)
		detailData := []warehouse.Instance{}
		json.Unmarshal([]byte(rec.Detail), &detailData)
		log.Printf("Quyen debug detail data %+v", detailData)
		ins.Detail = detailData
		dataArray = append(dataArray, ins)
	}
	ret.Data = dataArray

	return ret, nil
}

func (s *Storage) UpdateReceipt(in receipt.Instance, user string) (receipt.Instance, error) {
	record := Receipt{
		Name:   in.Name,
		Phone:  in.Phone,
		Adress: in.Addr,
		Status: in.Status,
	}
	result := db.Model(&record).Debug().Where(&Receipt{
		BaseModel: BaseModel{
			ID: in.ID,
		},
		User: user,
	}).Updates(&record)

	if result.Error != nil {
		return receipt.Instance{}, result.Error
	}

	return in, nil

}

func (s *Storage) DeleteReceipt(id int64, user string) error {
	temp, err := s.GetReceipt(uint64(id), user)
	if err != nil {
		return err
	}

	for _, item := range temp.Detail {
		ins, _ := s.GetItem(item.ID, user)
		ins.Remain = ins.Remain + int64(item.Quantity)
		s.UpdateItem(ins, user)
	}

	result := db.Model(&Receipt{}).Debug().Delete(&Receipt{
		BaseModel: BaseModel{
			ID: uint64(id),
		},
		User: user,
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *Storage) ReportReceipt(user string, from string, to string, status int) (receipt.Report, error) {
	var whs []Warehouse
	var recpts []Receipt

	timeFrom, err := time.Parse(time.RFC3339, from)
	if err != nil {
		return receipt.Report{}, err
	}

	timeTo, err := time.Parse(time.RFC3339, to)
	if err != nil {
		return receipt.Report{}, err
	}

	var result *gorm.DB

	if status == 255 {
		result = s.db.Debug().Model(&Receipt{}).
			Where("created_at >= ? AND created_at <= ? AND user = ?", timeFrom, timeTo, user).
			Find(&recpts)
		if result.Error != nil {
			return receipt.Report{}, nil
		}
	} else {
		result = s.db.Debug().Model(&Receipt{}).
			Where("created_at >= ? AND created_at <= ? AND user = ? AND status = ?", timeFrom, timeTo, user, status).
			Find(&recpts)
		if result.Error != nil {
			return receipt.Report{}, nil
		}
	}

	result = s.db.Debug().Model(&Warehouse{}).Where(&Warehouse{
		User: user,
	}).Find(&whs)
	if result.Error != nil {
		return receipt.Report{}, nil
	}

	var ret receipt.Report
	ret.TotalMn = 0
	ret.TotalProfit = 0

	for i := range whs {
		whs[i].Total = 0
		whs[i].Profit = 0
	}

	for _, rec := range recpts {
		log.Printf("Quyen debug detail data json: %+v", rec.Detail)
		detailData := []warehouse.Instance{}
		json.Unmarshal([]byte(rec.Detail), &detailData)
		log.Printf("Quyen debug detail data %+v", detailData)
		for _, dt := range detailData {
			for i := range whs {
				if dt.ID == whs[i].ID {
					whs[i].Total += int64(dt.Quantity) * int64(whs[i].SalePrice)
					whs[i].Profit += int64(dt.Quantity)*int64(whs[i].SalePrice) - int64(dt.Quantity)*int64(whs[i].PurchasePrice)
					whs[i].SaleTotal += dt.Quantity
					continue
				}
			}
		}
	}

	for _, item := range whs {
		ins := warehouse.Instance{
			ID:        item.ID,
			Name:      item.Name,
			SalePrice: item.SalePrice,
			Total:     item.Total,
			Profit:    item.Profit,
			SaleTotal: item.SaleTotal,
		}
		ret.Detail = append(ret.Detail, ins)
		ret.TotalMn += item.Total
		ret.TotalProfit += item.Profit
	}

	return ret, nil

}
