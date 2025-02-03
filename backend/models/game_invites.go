package models

import (
	"gorm.io/gorm"
)

type GameInvites struct {
	gorm.Model
	ReceiverID int  `json:"receiver_id"`
	SenderID   int  `json:"sender_id"`
	Sender     User `json:"sender" gorm:"foreignKey:SenderID;references:ID"`
	Receiver   User `json:"receiver" gorm:"foreignKey:ReceiverID;references:ID"`
}
