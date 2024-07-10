package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"log"
)

func main() {

	db = initDatabase()
	defer db.Close()
	createSqlObj()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// Маршрут для страницы регистрации
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Маршрут для страницы авторизации
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	// Маршрут для главной страницы
	r.GET("/main", func(c *gin.Context) {
		c.HTML(http.StatusOK, "main.html", nil)
	})

	// Маршрут для страницы чата
	r.GET("/chat/:chatName", func(c *gin.Context) {
		chatName := c.Param("chatName")
		c.HTML(http.StatusOK, "chat.html", gin.H{"name": chatName})
	})

	r.GET("/chat", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", nil)
	})

	// Маршрут
	r.POST("/main/createChat", createChatHandler)

	// Маршрут для регистрации
	r.POST("/register", registerHandler)

	// Маршрут для авторизации
	r.POST("/login", loginHandler)

	// Маршрут для выхода
	r.POST("/logout", logoutHandler)

	// Маршрут для аутентификации перед WebSocket соединением
	r.POST("/auth_ws/:chatUid", authWSHandler)

	// Маршрут для вебсокета
	r.GET("/ws/:chatUid", wsHandler)

	go handleBroadcast()

	// Удаляем протухшие чаты
	go func() {
		for {
			log.Println("Runnig deleteExpiredChats...")
			deleteExpiredChats()
			time.Sleep(1 * time.Minute)
		}
	}()

	log.Println("Server started on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
