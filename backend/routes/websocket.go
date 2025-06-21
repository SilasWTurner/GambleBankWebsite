package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/dgrijalva/jwt-go"
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
	fmt.Println("WebSocketHandler called")

	tokenString := c.Query("token")
	if tokenString == "" {
		fmt.Println("No token provided")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token is required"})
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		fmt.Println("Invalid token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	username := claims.Username
	if username == "" {
		fmt.Println("No username in token claims")
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	fmt.Println("Upgrading connection for user:", username)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Failed to upgrade connection:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upgrade connection"})
		return
	}

	fmt.Println("Connection upgraded for user:", username)
	connections.Lock()
	connections.m[username] = &WebSocketConnection{Conn: conn, User: username}
	connections.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				fmt.Println("Client disconnected:", username)
			} else {
				fmt.Println("Error reading message:", err)
			}
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

		fmt.Println("Message received:", string(message))
		var msg struct {
			Action string        `json:"action"`
			Args   []interface{} `json:"args"`
		}
		if err := json.Unmarshal(message, &msg); err != nil {
			fmt.Println("Invalid message format:", err)
			conn.WriteJSON(gin.H{"error": "invalid message format"})
			continue
		}

		switch msg.Action {
		case "accept_invite":
			fmt.Println("Accept invite msg received")
			if len(msg.Args) > 0 {
				fmt.Println(msg.Args[0])
				inviteIDFloat, ok := msg.Args[0].(float64) // JSON numbers are float64
				if ok {
					inviteID := int(inviteIDFloat)
					fmt.Println("Accepting invite:", inviteID)
					acceptInvite(username, inviteID)
				} else {
					fmt.Println("Invalid invite ID type")
				}
			}
		}
	}
}

func acceptInvite(username string, inviteID int) {
	fmt.Println("acceptInvite called with username:", username, "and inviteID:", inviteID)
	var invite models.GameInvites
	if err := database.DB.Where("id = ?", inviteID).First(&invite).Error; err != nil {
		fmt.Println("Invite not found:", err)
		return
	}

	var receiver models.User
	if err := database.DB.Where("username = ?", username).First(&receiver).Error; err != nil {
		fmt.Println("Receiver not found:", err)
		return
	}

	if invite.ReceiverID != int(receiver.ID) {
		fmt.Println("Invite receiver ID does not match")
		return
	}

	// Start the game logic here
	startGame(invite.SenderID, invite.ReceiverID)

	// Delete the invite after acceptance
	database.DB.Delete(&invite)
}

func startGame(p1ID, p2ID int) {
	fmt.Println("startGame called with player1ID:", p1ID, "and player2ID:", p2ID)
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
