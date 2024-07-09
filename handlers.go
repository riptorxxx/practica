package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	jwtKey    = []byte("my_secret_key")
	clients   = make(map[string]map[*websocket.Conn]bool) // clients[chatUid] = map[connections]
	broadcast = make(chan Message)
	upgrader  = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
)

type Claims struct {
	UserLogin string `json:"userlogin"`
	jwt.StandardClaims
}

type User struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type Chat struct {
	Name      string `json:"name"`
	Lifetime  string `json:"lifetime"`
	Cypher    string `json:"cypher"`
	UserToken string `json:"token"`
}

type Message struct {
	ChatUID  string `json:"chatUid"`
	Username string `json:"username"`
	Text     string `json:"text"`
}

// ____________________________________________  Маршрут для авторизации  ____________________________________________  \\
func loginHandler(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var dbUser User
	err := db.QueryRow("SELECT login, password FROM users WHERE login=$1", user.Login).Scan(&dbUser.Login, &dbUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Invalid login or password: %v", user)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password"})
		} else {
			log.Printf("Error querying user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if dbUser.Password != user.Password {
		log.Printf("Invalid password for user: %v", user.Login)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password"})
		return
	}

	expirationTime := time.Now().Add(60 * time.Minute) // протухание jwt token
	claims := &Claims{
		UserLogin: user.Login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error creating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	// Сохранение токена в базе данных
	_, err = db.Exec("UPDATE users SET token=$1 WHERE login=$2", tokenString, user.Login)
	if err != nil {
		log.Printf("Error updating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	log.Printf("User logged in successfully: %v", user.Login)
	// c.JSON(http.StatusOK, gin.H{"token": tokenString})
	c.JSON(http.StatusOK, gin.H{"token": tokenString, "redirect": "/main"})
}

// ____________________________________________  Маршрут для регистрации  ____________________________________________  \\
func registerHandler(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Проверка уникальности
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE login=$1 OR email=$2 OR phone=$3)",
		user.Login, user.Email, user.Phone).Scan(&exists)
	if err != nil {
		log.Printf("Error checking uniqueness: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists {
		log.Printf("User with login, email or phone already exists: %v", user)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Login, Email or Phone already exists"})
		return
	}

	// Сохранение пользователя
	_, err = db.Exec("INSERT INTO users (login, email, phone, password) VALUES ($1, $2, $3, $4)",
		user.Login, user.Email, user.Phone, user.Password)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	log.Printf("User registered successfully: %v %s %s %s", user, " Email: "+user.Email, " Phone: "+user.Phone, " Pass: "+user.Password)
	// c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "redirect": "/login"})
}

// ____________________________________________  Маршрут для выхода  ____________________________________________  \\
func logoutHandler(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		log.Printf("Token required for logout")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		log.Printf("Invalid token for logout: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	_, err = db.Exec("UPDATE users SET token='' WHERE login=$1", claims.UserLogin)
	if err != nil {
		log.Printf("Error clearing token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	log.Printf("User logged out successfully: %v", claims.UserLogin)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// ____________________________________________  Auth Websocket Handler  ____________________________________________  \\
func authWSHandler(c *gin.Context) {
	chatUid := c.Param("chatUid")
	tokenString := c.PostForm("token")

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		return
	}

	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !parsedToken.Valid {
		log.Printf("Invalid token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.SetCookie("chatUid", chatUid, 3600, "/", "", false, true)
	c.SetCookie("token", tokenString, 3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Authenticated", "username": claims.UserLogin})
}

// ____________________________________________  Websocket Handler  ____________________________________________  \\

func wsHandler(c *gin.Context) {
	chatUid, err := c.Cookie("chatUid")
	if err != nil {
		log.Printf("Chat UID cookie error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Chat UID is required"})
		return
	}

	tokenString, err := c.Cookie("token")
	if err != nil {
		log.Printf("Token cookie error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		log.Printf("Invalid token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if clients[chatUid] == nil {
		clients[chatUid] = make(map[*websocket.Conn]bool)
	}
	clients[chatUid][conn] = true
	log.Printf("WebSocket connection established for chat %v: %v", chatUid, conn.RemoteAddr())

	go handleWebSocketConnection(conn, chatUid, claims.UserLogin)
}

func handleWebSocketConnection(conn *websocket.Conn, chatUid, username string) {
	defer func() {
		conn.Close()
		delete(clients[chatUid], conn)
		if len(clients[chatUid]) == 0 {
			delete(clients, chatUid)
		}
		log.Printf("WebSocket connection closed for chat %v: %v", chatUid, conn.RemoteAddr())
	}()

	for {
		var message Message
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("Error reading JSON from WebSocket: %v", err)
			break
		}

		message.Username = username
		message.ChatUID = chatUid
		log.Printf("Received message from %v in chat %v: %v", message.Username, message.ChatUID, message.Text)
		broadcast <- message
	}
}

func handleBroadcast() {
	for {
		message := <-broadcast

		// Определяем комнату чата из сообщения
		chatUid := message.ChatUID

		for client := range clients[chatUid] {
			err := client.WriteJSON(message)
			if err != nil {
				log.Printf("Ошибка записи JSON в WebSocket: %v", err)
				client.Close()
				delete(clients[chatUid], client)
			}
		}
	}
}

// ____________________________________________  Create Chat Handler  ____________________________________________  \\
func createChatHandler(c *gin.Context) {
	var chat Chat
	if err := c.BindJSON(&chat); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(chat.UserToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		log.Printf("Invalid token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var userUID string
	err = db.QueryRow("SELECT uid FROM users WHERE login=$1", claims.UserLogin).Scan(&userUID)
	if err != nil {
		log.Printf("Error fetching user UID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	lifetime, err := time.ParseDuration(chat.Lifetime)
	if err != nil {
		log.Printf("Error parsing lifetime: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lifetime format"})
		return
	}

	expirationTime := time.Now().Add(lifetime)

	_, err = db.Exec("INSERT INTO chats (name, lcc, cypher, user_uid) VALUES ($1, $2, $3, $4)",
		chat.Name, expirationTime, chat.Cypher, userUID)
	if err != nil {
		log.Printf("Error creating chat: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	log.Printf("Chat created successfully: %v", chat.Name)
	c.JSON(http.StatusCreated, gin.H{"message": "Chat created successfully", "redirect": "/main"})
}
