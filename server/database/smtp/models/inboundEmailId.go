package models

import "gorm.io/gorm"

type InboundEmailId struct {
	gorm.Model
	EmailId string `gorm:"uniqueIndex"`
}

func (ie InboundEmailId) Marshal() (string, error) {
	return ie.EmailId, nil
}

func UnmarshalInboundEmailId(data string) (InboundEmailId, error) {
	return InboundEmailId{EmailId: data}, nil
}
