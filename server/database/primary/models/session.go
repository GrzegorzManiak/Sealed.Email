package models

import "gorm.io/gorm"

type Session struct {
	gorm.Model
	User   User `gorm:"foreignKey:UserID" json:"user"`
	UserID uint `gorm:"column:uid"        json:"uid"`

	SessionID string `json:"sid"            orm:"column:sid;uniqueIndex"`
	ExpiresAt int64  `gorm:"column:exp"     json:"exp"`
	CreatedAt int64  `gorm:"column:iat"     json:"iat"`
	RefreshAt int64  `gorm:"column:rat"     json:"rat"`
	Group     string `gorm:"column:grp"     json:"grp"`
	Revoked   bool   `gorm:"column:revoked" json:"revoked"`
}
