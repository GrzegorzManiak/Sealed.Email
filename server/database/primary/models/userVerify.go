package models

import "gorm.io/gorm"

type UserVerify struct {
	gorm.Model
	RID    string `gorm:"uniqueIndex"`
	UserID uint

	// -- ClientAuthInitRequestPayload
	X1    string
	X2    string
	PI1_V string
	PI1_R string
	PI2_V string
	PI2_R string

	// -- ServerAuthInitResponse
	XPub4    string
	XPri4    string
	Beta     string
	PI4_V    string
	PI4_R    string
	PIBeta_V string
	PIBeta_R string
}
