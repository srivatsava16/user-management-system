package models

import (
	"time"

	"gorm.io/datatypes"
)

type AuditLog struct {
	ID            int            `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	CreatedAt     time.Time      `json:"created_at" gorm:"column:created_at"`
	ChangeBy      string         `json:"change_by" gorm:"column:changed_by"`
	ChangeType    string         `json:"change_type" gorm:"column:change_type"`
	ChangeDetails string         `json:"change_details" gorm:"column:change_details"`
	BeforeJson    datatypes.JSON `json:"before_json" gorm:"column:before_json;type:json"`
	AfterJson     datatypes.JSON `json:"after_json" gorm:"column:after_json;type:json"`
}

func (AuditLog) TableName() string {
	return "APT_CUSTOM_ZXDS_AUDIT_LOGS"
}
