package models

import "gorm.io/gorm"

type Session struct {
	gorm.Model
	User   User `json:"user" gorm:"foreignKey:UserID"`
	UserID uint `json:"uid" gorm:"column:uid"`

	SessionID string `json:"sid" orm:"column:sid;uniqueIndex"`
	ExpiresAt int64  `json:"exp" gorm:"column:exp"`
	CreatedAt int64  `json:"iat" gorm:"column:iat"`
	RefreshAt int64  `json:"rat" gorm:"column:rat"`
	Group     string `json:"grp" gorm:"column:grp"`
	Revoked   bool   `json:"revoked" gorm:"column:revoked"`
}
