package gormx

import (
	"database/sql/driver"
	"encoding/json"
)

type Model struct {
	ID        int64 `json:"id" gorm:"primarykey,column:id" form:"id"`
	CreatedAt int64 `json:"created_at" gorm:"column:created_at" form:"created_at"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at" form:"updated_at"`
}

type Dict map[string]interface{}

func NewDict(v interface{}) Dict {
	d := make(Dict)
	marshal, err := json.Marshal(v)
	if err != nil {
		return d
	}
	_ = json.Unmarshal(marshal, &d)
	return d
}

func (d Dict) Value() (driver.Value, error) {
	marshal, err := json.Marshal(d)
	if err != nil {
		return "{}", nil
	}
	return string(marshal), nil
}

func (d *Dict) Scan(v interface{}) error {
	bytes := v.([]byte)
	err := json.Unmarshal(bytes, d)
	if err != nil {
		*d = make(Dict)
	}
	return nil
}
