package models

import "time"

type BusinessUnit struct {
	ID          int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"column:name;size:100;uniqueIndex"`
	Description string    `json:"description" gorm:"column:description"`
	Code        string    `json:"code" gorm:"column:code;size:50;uniqueIndex"`
	IsActive    *bool     `json:"is_active" gorm:"column:is_active"`
	ApprovedBy  *string   `json:"approved_by,omitempty" gorm:"column:approved_by"`
	CreatedBy   *string   `json:"created_by,omitempty" gorm:"column:created_by"`
	UpdatedBy   *string   `json:"updated_by,omitempty" gorm:"column:updated_by"`
	CreatedAt   time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
}

func (BusinessUnit) TableName() string {
	return "APT_CUSTOM_ZXDS_BUSINESS_UNITS"
}
