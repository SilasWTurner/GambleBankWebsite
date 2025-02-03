package models

import (
	"gorm.io/gorm"
)

type CompletedGames struct {
	gorm.Model
	Player1ID    int  `json:"player1_id"`
	Player2ID    int  `json:"player2_id"`
	WinnerID     int  `json:"winner_id"`
	Player1      User `json:"player1" gorm:"foreignKey:Player1ID;references:ID"`
	Player2      User `json:"player2" gorm:"foreignKey:Player2ID;references:ID"`
	Winner       User `json:"winner" gorm:"foreignKey:WinnerID;references:ID"`
	Player1Score int  `json:"player1_score"`
	Player2Score int  `json:"player2_score"`
}
