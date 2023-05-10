package dao

import (
	"time"

	"gorm.io/gorm"

	"github.com/mapprotocol/brc20Relayer/resource/db"
)

type Test struct {
	ID        uint64         `gorm:"primarykey"`
	Info      string         `gorm:"column:info" json:"info" sql:"json"`
	TxHash    string         `gorm:"column:tx_hash" json:"tx_hash" sql:"char(66)"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at" sql:"datetime"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at" sql:"datetime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at" sql:"datetime"`
}

func (t *Test) TableName() string {
	return "test"
}

func (t *Test) Create() (id uint64, err error) {
	err = db.GetDB().Create(t).Error
	return t.ID, err
}

func (t *Test) Get() (test *Test, err error) {
	err = db.GetDB().Where(t).First(&test).Error
	return test, err
}
