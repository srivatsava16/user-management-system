package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Role struct {
	ID             int                 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name           string              `json:"name" gorm:"column:name;size:255;uniqueIndex:idx_role_name_div"`
	RoleType       string              `json:"role_type" gorm:"column:role_type"`
	Description    string              `json:"description" gorm:"column:description"`
	BusinessUnitID int                 `json:"business_unit_id" gorm:"column:business_unit_id"`
	DivisionID     int                 `json:"division_id" gorm:"column:division_id;uniqueIndex:idx_role_name_div"`
	IsActive       *bool               `json:"is_active" gorm:"column:is_active"`
	DataSetAccess  JSON                `json:"data_set_access" gorm:"column:data_set_access;type:json"`
	Permissions    []ModulePermissions `json:"permissions" gorm:"foreignKey:RoleID"`
	CreatedBy      *string             `json:"created_by,omitempty" gorm:"column:created_by"`
	UpdatedBy      *string             `json:"updated_by,omitempty" gorm:"column:updated_by"`
	CreatedAt      time.Time           `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt      *time.Time          `json:"updated_at,omitempty" gorm:"column:updated_at"`
}

func (Role) TableName() string {
	return "APT_CUSTOM_ZXDS_ROLES"
}

type ModulePermissions struct {
	ID                  int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	RoleID              int    `json:"role_id" gorm:"column:role_id"`
	ModuleName          string `json:"module_name" gorm:"column:module_name"`
	WithinDivision      JSON   `json:"within_division" gorm:"column:within_division;type:json"`
	AcrossDivisions     JSON   `json:"across_divisions" gorm:"column:across_divisions;type:json"`
	WithinBusinessUnit  JSON   `json:"within_business_unit" gorm:"column:within_business_unit;type:json"`
	AcrossBusinessUnits JSON   `json:"across_business_units" gorm:"column:across_business_units;type:json"`
	SpecialPermissions  JSON   `json:"special_permissions" gorm:"column:special_permissions;type:json"`
}

func (ModulePermissions) TableName() string {
	return "APT_CUSTOM_ZXDS_MODULE_PERMISSIONS"
}

type JSON map[string]interface{}

func (j JSON) Value() (driver.Value, error) {
	bytes, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}
func (j *JSON) Scan(value interface{}) error {
	// bytes, ok := value.([]byte)
	// if !ok {
	// 	return nil
	// }
	if value == nil {
		*j = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("JSON.Scan: unsupported type: %T, value: %v", value, value)
	}
	return json.Unmarshal(bytes, j)
}
