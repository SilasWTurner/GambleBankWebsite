package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silaswturner/GambleBankWebsite/backend/database"
	"github.com/silaswturner/GambleBankWebsite/backend/models"
)

type GameInviteResponse struct {
	ID         uint   `json:"id"`
	SenderID   uint   `json:"sender_id"`
	SenderName string `json:"sender_name"`
	CreatedAt  string `json:"created_at"`
}

func SendInvite(c *gin.Context) {
	// Get the sender's username from the context
	senderUsername, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "username not found in context"})
		return
	}

	// Get the sender's user information from the database
	var sender models.User
	if err := database.DB.Where("username = ?", senderUsername).First(&sender).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "sender not found"})
		return
	}

	// Get the receiver's username from the request body
	var requestBody struct {
		ReceiverUsername string `json:"receiver_username"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the receiver's username and get their user information from the database
	var receiver models.User
	if err := database.DB.Where("username = ?", requestBody.ReceiverUsername).First(&receiver).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "receiver not found"})
		return
	}

	// Check if an invite already exists
	var existingInvite models.GameInvites
	if err := database.DB.Where("sender_id = ? AND receiver_id = ?", sender.ID, receiver.ID).First(&existingInvite).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invite already exists"})
		return
	}

	// Add an entry in the game_invites table
	gameInvite := models.GameInvites{
		SenderID:   int(sender.ID),
		ReceiverID: int(receiver.ID),
	}
	if err := database.DB.Create(&gameInvite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create game invite"})
		return
	}

	// Broadcast the invite to the receiver if they are connected
	connections.RLock()
	if conn, ok := connections.m[receiver.Username]; ok {
		conn.Conn.WriteJSON(gin.H{"action": "new_invite", "args": []interface{}{gameInvite.ID, sender.Username}})
	}
	connections.RUnlock()

	c.JSON(http.StatusOK, gin.H{"message": "game invite sent successfully"})
}

func ListInvites(c *gin.Context) {
	// Get the receiver's username from the context
	receiverUsername, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "username not found in context"})
		return
	}

	// Get the receiver's user information from the database
	var receiver models.User
	if err := database.DB.Where("username = ?", receiverUsername).First(&receiver).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "receiver not found"})
		return
	}

	// Get all game invites for the receiver and join with the users table to get the sender's username
	var gameInvites []GameInviteResponse
	if err := database.DB.Table("game_invites").
		Select("game_invites.id, game_invites.sender_id, users.username as sender_name, game_invites.created_at").
		Joins("left join users on users.id = game_invites.sender_id").
		Where("game_invites.receiver_id = ?", receiver.ID).
		Scan(&gameInvites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve game invites"})
		return
	}

	c.JSON(http.StatusOK, gameInvites)
}

func AcceptInvite(c *gin.Context) {
	// Get the receiver's username from the context
	receiverUsername, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "username not found in context"})
		return
	}

	// Get the invite ID from the request body
	var requestBody struct {
		InviteID int `json:"invite_id"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acceptInvite(receiverUsername.(string), requestBody.InviteID)

	c.JSON(http.StatusOK, gin.H{"message": "invite accepted"})
}
