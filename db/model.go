package db

import (
	"database/sql/driver"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type LocalTime time.Time

type BaseModel struct {
	ID        uint           `gorm:"primarykey;autoIncrement" json:"id"`
	CreatedAt *LocalTime     `json:"created_at"`
	UpdatedAt *LocalTime     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (d *LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*d)
	return fmt.Appendf(nil, "\"%v\"", tTime.Format(time.DateTime)), nil
}

func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(t)
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}

func (t *LocalTime) Scan(v any) error {
	if value, ok := v.(time.Time); ok {
		*t = LocalTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
