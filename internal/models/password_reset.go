package models

import (
	"time"
)

type PasswordResetToken struct {
	ID        int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Email     string    `json:"email" gorm:"column:email"`
	Token     string    `json:"token" gorm:"column:token;unique"`
	ExpiresAt time.Time `json:"expires_at" gorm:"column:expires_at"`
	Used      bool      `json:"used" gorm:"column:used;default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

func (PasswordResetToken) TableName() string {
	return "APT_CUSTOM_ZXDS_PASSWORD_RESET_TOKENS"
}
