package models

import (
	"time"
)

type TokenInvalidation struct {
	ID        int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID    int       `json:"user_id" gorm:"column:user_id"`
	Reason    string    `json:"reason" gorm:"column:reason"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	JTI       string    `json:"jti" gorm:"column:jti"`
}

func (TokenInvalidation) TableName() string {
	return "APT_CUSTOM_ZXDS_TOKEN_INVALIDATIONS"
}
