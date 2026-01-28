package models

import "time"

type Division struct {
	ID             int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name           string    `json:"name" gorm:"column:name;size:100;uniqueIndex:idx_division_name_bu"`
	Description    string    `json:"description" gorm:"column:description"`
	BusinessUnitID int       `json:"business_unit_id" gorm:"column:business_unit_id;uniqueIndex:idx_division_name_bu"`
	ApprovedBy     *string   `json:"approved_by,omitempty" gorm:"column:approved_by"`
	IsActive       *bool     `json:"is_active" gorm:"column:is_active;default:false"`
	CreatedBy      *string   `json:"created_by,omitempty" gorm:"column:created_by"`
	UpdatedBy      *string   `json:"updated_by,omitempty" gorm:"column:updated_by"`
	CreatedAt      time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt      time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
}

func (Division) TableName() string {
	return "APT_CUSTOM_ZXDS_DIVISIONS"
}
