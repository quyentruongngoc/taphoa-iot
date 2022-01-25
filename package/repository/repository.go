package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type BaseModel struct {
	ID        uint64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Config struct {
	Host   string
	Port   string
	Name   string
	User   string
	Passwd string
	Type   string
	Debug  bool
}

type Storage struct {
	db *gorm.DB
}

var initialized bool = false
var db *gorm.DB

func NewStorage() (*Storage, error) {
	return &Storage{db}, nil
}

func InitStorage(c Config) error {
	var err error
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=25s",
		c.User,
		c.Passwd,
		c.Host,
		c.Port,
		c.Name,
	)

	log.Println("Init storage with connection string:", connStr)

	db, err = gorm.Open(c.Type, connStr)
	if err != nil {
		err := errors.Wrapf(err, "Failed to open connection to database, connStr: %s", connStr)
		log.Println("Error", err)
		return err
	}

	if err = db.DB().Ping(); err != nil {
		err := errors.Wrapf(err, "Failed to connect to database %s", c.Name)
		log.Fatalln("Error", err)
		return err
	}

	db.SingularTable(true)
	if c.Debug {
		db.LogMode(true)
	}

	var models []interface{}
	models = []interface{}{
		&Account{},
		&Warehouse{},
		&Receipt{},
		&Customer{},
	}
	db.AutoMigrate(models...)

	log.Println("Seeding data to table account")
	for _, acc := range accountDf {
		db.FirstOrCreate(&acc, &acc)
	}

	return nil
}
