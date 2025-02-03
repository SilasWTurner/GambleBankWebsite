package routes

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/silaswturner/GambleBankWebsite/backend/database"
	"github.com/silaswturner/GambleBankWebsite/backend/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketConnection struct {
	Conn *websocket.Conn
	User string
}

var connections = struct {
	sync.RWMutex
	m map[string]*WebSocketConnection
}{m: make(map[string]*WebSocketConnection)}

func WebSocketHandler(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upgrade connection"})
		return
	}

	connections.Lock()
	connections.m[username] = &WebSocketConnection{Conn: conn, User: username}
	connections.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			connections.Lock()
			delete(connections.m, username)
			connections.Unlock()

			// Delete all invites from the database when the sender disconnects
			var sender models.User
			if err := database.DB.Where("username = ?", username).First(&sender).Error; err == nil {
				database.DB.Where("sender_id = ?", sender.ID).Delete(&models.GameInvites{})
			}

			conn.Close()
			break
		}

		var msg struct {
			Action string        `json:"action"`
			Args   []interface{} `json:"args"`
		}
		if err := json.Unmarshal(message, &msg); err != nil {
			conn.WriteJSON(gin.H{"error": "invalid message format"})
			continue
		}

		switch msg.Action {
		case "accept_invite":
			if len(msg.Args) > 0 {
				inviteID, ok := msg.Args[0].(float64) // JSON numbers are float64
				if ok {
					acceptInvite(username, int(inviteID))
				}
			}
		}
	}
}

func acceptInvite(username string, inviteID int) {
	var invite models.GameInvites
	if err := database.DB.Where("id = ?", inviteID).First(&invite).Error; err != nil {
		return
	}

	var receiver models.User
	if err := database.DB.Where("username = ?", username).First(&receiver).Error; err != nil {
		return
	}

	if invite.ReceiverID != int(receiver.ID) {
		return
	}

	// Start the game logic here
	startGame(invite.SenderID, invite.ReceiverID)

	// Delete the invite after acceptance
	database.DB.Delete(&invite)
}

func startGame(p1ID, p2ID int) {
	// Implement game logic here
	// For example, create a new game entry in the database
	game := models.OngoingGames{
		Player1ID: p1ID,
		Player2ID: p2ID,
	}
	database.DB.Create(&game)

	// Notify both players about the game start
	connections.RLock()
	if conn, ok := connections.m[getUsernameByID(p1ID)]; ok {
		conn.Conn.WriteJSON(gin.H{"action": "game_started", "args": []interface{}{getUsernameByID(p2ID)}})
	}
	if conn, ok := connections.m[getUsernameByID(p2ID)]; ok {
		conn.Conn.WriteJSON(gin.H{"action": "game_started", "args": []interface{}{getUsernameByID(p1ID)}})
	}
	connections.RUnlock()
}

func getUsernameByID(userID int) string {
	var user models.User
	database.DB.Where("id = ?", userID).First(&user)
	return user.Username
}
