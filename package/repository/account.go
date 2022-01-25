package repository

import (
	"fmt"
	"taphoa-iot-backend/package/api-interface/account"
	"taphoa-iot-backend/package/constant"
)

type Account struct {
	BaseModel

	User    string `sql:"type:varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci;unique;not null"`
	Passwd  string `sql:"type:varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci;not null"`
	Role    int
	Name    string `sql:"type:varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci"`
	Addr    string `sql:"type:text CHARACTER SET utf8 COLLATE utf8_general_ci"`
	Phone   string `sql:"type:varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci"`
	Email   string `sql:"type:varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci"`
	Gender  bool
	Shop    string `sql:"type:varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci;not null"`
	License int
}

var accountDf = []Account{
	{User: "admin", Passwd: "15da766985313b89b9138027b19131adf8632d061344f0209fa6af20e6808f71", Role: constant.AdminRole},
	{User: "nhung", Passwd: "d423e019ab0f44e422da0cf6ec5dfacda9fe24606993eccca59acc59671887a0", Role: constant.OwnerRole},
}

func (s *Storage) CreateUser(in account.Instance) error {
	// if isCreate {
	// 	checkAccount, _ := s.DescribeUser(in.Account.User)
	// 	if checkAccount.Account.User == in.Account.User {
	// 		return fmt.Errorf("Duplicate")
	// 	}
	// }

	// user := Account{
	// 	User:         in.Account.User,
	// 	DoctorUUID:   in.MedicalInfo.Doctor.User,
	// 	ExpertUUID:   in.MedicalInfo.Expert.User,
	// 	Passwd:       in.Account.Passwd,
	// 	Role:         in.Account.Role,
	// 	Name:         in.Mgmt.Name,
	// 	Addr:         in.Mgmt.Addr,
	// 	Phone:        in.Mgmt.Phone,
	// 	Email:        in.Mgmt.Email,
	// 	IDCard:       in.Mgmt.IDcard,
	// 	Gender:       in.Mgmt.Gender,
	// 	RelPhone:     in.Mgmt.RelPhone,
	// 	Sevirity:     0,
	// 	Discharge:    in.Mgmt.Discharge,
	// 	Subclinical:  in.MedicalInfo.Subclinical,
	// 	SelfHistory:  in.MedicalInfo.SelfHistory,
	// 	Vaccine:      in.MedicalInfo.Vaccine,
	// 	Birthday:     in.Mgmt.Birthday,
	// 	ReceivedTime: in.Mgmt.ReceivedTime,
	// 	DiseaseDate:  in.MedicalInfo.DiseaseDate,
	// }

	// if isCreate {
	// 	// log.Printf("Quyen debug time birthday: %v\n", in.Mgmt.Birthday)
	// 	// t, err := time.Parse(time.RFC3339, in.Mgmt.Birthday)
	// 	// if err != nil {
	// 	// 	user.Birthday = time.Now()
	// 	// 	log.Printf("failed to parse time format: %v - %v", in.Mgmt.Birthday, err)
	// 	// 	// return patient.PatientMgmt{}, fmt.Errorf("failed to parse time format: %v - %v", in.Transfer.Time, err)
	// 	// }
	// 	// user.Birthday = t
	// 	// log.Printf("Quyen debug time birthday parse: %v\n", user.Birthday.Format(time.RFC3339))

	// 	// t, err = time.Parse(time.RFC3339, in.Mgmt.ReceivedTime)
	// 	// if err != nil {
	// 	// 	user.ReceivedTime = time.Now()
	// 	// 	// return patient.PatientMgmt{}, fmt.Errorf("failed to parse time format: %v - %v", in.Transfer.Time, err)
	// 	// }
	// 	// user.ReceivedTime = t
	// 	// log.Printf("Quyen debug time recervied time parse: %v\n", user.ReceivedTime.Format(time.RFC3339))

	// 	// t, err = time.Parse(time.RFC3339, in.MedicalInfo.DiseaseDate)
	// 	// if err != nil {
	// 	// 	user.DiseaseDate = time.Now()
	// 	// 	// return patient.PatientMgmt{}, fmt.Errorf("failed to parse time format: %v - %v", in.Transfer.Time, err)
	// 	// }
	// 	// user.DiseaseDate = t
	// 	// log.Printf("Quyen debug time DiseaseDate parse: %v\n", user.DiseaseDate.Format(time.RFC3339))

	// 	// log.Printf("Quyen debug time user parse: %+v\n", user)
	// 	var err error

	// 	if user.Role == constant.PatientRole {
	// 		if user.DoctorUUID == "" {
	// 			user.DoctorUUID, err = s.getDoctorUUID()
	// 			if err != nil {
	// 				return err
	// 			}
	// 		}
	// 	}

	// 	if user.Role == constant.DoctorRole || user.Role == constant.AutoDoctorRole {
	// 		if user.ExpertUUID == "" {
	// 			user.ExpertUUID, err = s.getExpertUUID()
	// 			if err != nil {
	// 				return err
	// 			}
	// 		}
	// 	}

	// 	result := s.db.Debug().FirstOrCreate(&user, &user)
	// 	if result.Error != nil {
	// 		return fmt.Errorf("Duplicate")
	// 	}
	// } else {
	// 	result := db.Model(&user).Debug().Where(&Account{
	// 		User: in.Account.User,
	// 	}).Updates(&user)

	// 	if result.Error != nil {
	// 		return result.Error
	// 	}
	// }

	return nil
}

func (s *Storage) DescribeUser(user string) (account.Instance, error) {
	var record Account
	var ret account.Instance

	result := s.db.Debug().Where(&Account{
		User: user,
	}).First(&record)
	if result.Error != nil {
		return account.Instance{}, fmt.Errorf("Failed to search user in database: %+v : %+v", user, result.Error)
	}

	ret = account.Instance{
		User:   record.User,
		Passwd: record.Passwd,
		Role:   record.Role,
		Token:  "",
		Name:   record.Name,
		Addr:   record.Addr,
		Phone:  record.Phone,
		Email:  record.Email,
		Gender: record.Gender,
	}

	return ret, nil
}
