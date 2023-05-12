package dao

import (
	"time"

	"gorm.io/gorm"

	"github.com/mapprotocol/brc20Relayer/resource/db"
)

type BRC20 struct {
	ID            uint64         `gorm:"primarykey"`
	Metadata      string         `gorm:"column:Metadata" json:"Metadata" sql:"json"`
	Height        uint64         `gorm:"column:height" json:"height" sql:"bigint(20)"`
	InscriptionID string         `gorm:"column:inscription_id" json:"inscription_id" sql:"char(66)"`
	TxID          string         `gorm:"column:tx_id" json:"tx_id" sql:"char(64)"`
	TxIdx         uint64         `gorm:"column:tx_idx" json:"tx_idx" sql:"bigint(20)"`
	Account       string         `gorm:"column:account" json:"account" sql:"char(42)"`
	TxHash        string         `gorm:"column:tx_hash" json:"tx_hash" sql:"char(66)"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at" sql:"datetime"`
	UpdatedAt     time.Time      `gorm:"column:updated_at" json:"updated_at" sql:"datetime"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at" sql:"datetime"`
}

func (b *BRC20) TableName() string {
	return "brc20"
}

func NewBRC20() *BRC20 {
	return new(BRC20)
}

func (b *BRC20) Create() (id uint64, err error) {
	err = db.GetDB().Create(b).Error
	return b.ID, err
}

func (b *BRC20) BatchCreate(bs []*BRC20) error {
	return db.GetDB().Create(bs).Error
}

func (b *BRC20) Get() (b20 *BRC20, err error) {
	err = db.GetDB().Where(b).First(&b20).Error
	return b20, err
}
