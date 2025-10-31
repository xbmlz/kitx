package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type DateTime time.Time

type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt DateTime       `json:"created_at"`
	UpdatedAt DateTime       `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (d *DateTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*d)
	return fmt.Appendf(nil, "\"%v\"", tTime.Format(time.DateTime)), nil
}
