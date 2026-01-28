package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type User struct {
	ID             int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	FullName       string    `json:"full_name" gorm:"column:full_name;size:100;uniqueIndex"`
	Email          string    `json:"email" gorm:"column:email;size:100;uniqueIndex"`
	Password       string    `json:"password" gorm:"column:password"`
	UserName       string    `json:"user_name" gorm:"column:user_name;size:100;uniqueIndex"`
	BusinessUnitID int       `json:"business_unit_id" gorm:"column:business_unit_id"`
	DivisionID     int       `json:"division_id" gorm:"column:division_id"`
	IsActive       *bool     `json:"is_active" gorm:"column:is_active;default:false"`
	RoleIds        IntArray  `json:"role_ids" gorm:"column:role_ids;type:json"`
	CreatedBy      *string   `json:"created_by,omitempty" gorm:"column:created_by"`
	UpdatedBy      *string   `json:"updated_by,omitempty" gorm:"column:updated_by"`
	ApprovedBy     *string   `json:"approved_by,omitempty" gorm:"column:approved_by"`
	CreatedAt      time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt      time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "APT_CUSTOM_ZXDS_USERS"
}

type UserPermisssionsResponse struct {
	UserID    int    `json:"user_id"`
	UserEmail string `json:"user_email"`
	Roles     []Role `json:"roles"`
}

type IntArray []int

func (j *IntArray) Scan(value interface{}) error {
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
		return errors.New("failed to scan IntArray: incompatible type")
	}
	return json.Unmarshal(bytes, j)
}

func (j IntArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}

	bytes, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}

	return string(bytes), nil
}
