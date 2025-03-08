package db

import (
	"github.com/koyo-os/tokino/src/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Data struct{
	db *gorm.DB
}

func New() (*Data, error) {
	db,err := gorm.Open(sqlite.Open("../storage/main.db"))
	if err != nil{
		return nil,err
	}
	
	err = db.AutoMigrate(&models.Block{})
	if err != nil{
		return nil,err
	}

	return &Data{
		db: db,
	}, nil
}

func (d *Data) Add(b *models.Block) error {
	return d.db.Create(b).Error
}

func (d *Data) GetAll() ([]models.Block, error) {
	var blocks []models.Block

	res := d.db.Find(&blocks)
	return blocks, res.Error
}

func (d *Data) GetLast() (models.Block, error) {
	var block models.Block

	res := d.db.Last(&block)
	return block, res.Error
}