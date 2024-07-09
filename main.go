package main

import (
	"net/http"
	"time"

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

// Функции для работы с WebSocket
// func handleWebSocketConnection(conn *websocket.Conn, chatUid string) {
// 	defer func() {
// 		conn.Close()
// 		delete(clients[chatUid], conn)
// 		log.Printf("WebSocket connection closed: %v", conn.RemoteAddr())
// 	}()

// 	for {
// 		var message Message
// 		err := conn.ReadJSON(&message)
// 		if err != nil {
// 			log.Printf("Error reading JSON from WebSocket: %v", err)
// 			break
// 		}

// 		log.Printf("Received message from %v in chat %v: %v", message.Username, chatUid, message.Text)
// 		broadcast <- message
// 	}
// }

// ________________________________________________          ______________________________________________________

// func handleWebSocketConnection(conn *websocket.Conn, chatUid string) {
// 	defer func() {
// 		conn.Close()
// 		delete(clients[chatUid], conn)
// 		log.Printf("WebSocket connection closed: %v", conn.RemoteAddr())
// 	}()

// 	for {
// 		var message Message
// 		err := conn.ReadJSON(&message)
// 		if err != nil {
// 			log.Printf("Error reading JSON from WebSocket: %v", err)
// 			break
// 		}

// 		log.Printf("Received message from %v in chat %v: %v", message.Username, chatUid, message.Text)
// 		broadcast <- message
// 	}
// }

// func handleBroadcast() {
// 	for {
// 		message := <-broadcast

// 		// Определяем комнату чата из сообщения
// 		chatUid := message.ChatUID

// 		for client := range clients[chatUid] {
// 			err := client.WriteJSON(message)
// 			if err != nil {
// 				log.Printf("Ошибка записи JSON в WebSocket: %v", err)
// 				client.Close()
// 				delete(clients[chatUid], client)
// 			}
// 		}
// 	}
// }

// func handleWebSocketConnection(conn *websocket.Conn, chatUid string) {
// 	defer func() {
// 		conn.Close()
// 		delete(clients[chatUid], conn)
// 		log.Printf("WebSocket connection closed: %v", conn.RemoteAddr())
// 	}()

// 	for {
// 		var message Message
// 		err := conn.ReadJSON(&message)
// 		if err != nil {
// 			log.Printf("Error reading JSON from WebSocket: %v", err)
// 			break
// 		}

// 		log.Printf("Received message from %v: %v", message.Username, message.Text)
// 		broadcast <- message
// 	}
// }

// func handleBroadcast() {
// 	for {
// 		message := <-broadcast

// 		// Определяем комнату чата из сообщения или используем общий чат, если нужно
// 		chatRoom := "" // Извлекаем комнату чата из сообщения или используем общий чат, если применимо

// 		for client := range clients {
// 			// Проверяем, находится ли клиент в той же комнате чата, что и отправитель сообщения
// 			// Подстраиваем эту логику в зависимости от управления комнатами чата
// 			if isInChatRoom(client, chatRoom) {
// 				err := client.WriteJSON(message)
// 				if err != nil {
// 					log.Printf("Ошибка записи JSON в WebSocket: %v", err)
// 					client.Close()
// 					delete(clients, client)
// 				}
// 			}
// 		}
// 	}
// }

// // Пример функции для проверки, находится ли клиент в той же комнате чата
// func isInChatRoom(client *websocket.Conn, chatRoom string) bool {
// 	// Реализуйте логику для определения, находится ли клиент в указанной комнате чата
// 	// Пример: Сравните комнату чата клиента с указанной комнатой чата
// 	return true // Заглушка; замените на реальную логику
// }

// func handleBroadcast() {
// 	for {
// 		message := <-broadcast

// 		for client := range clients {
// 			err := client.WriteJSON(message)
// 			if err != nil {
// 				log.Printf("Error writing JSON to WebSocket: %v", err)
// 				client.Close()
// 				delete(clients, client)
// 			}
// 		}
// 	}
// }
// Обновленная функция handleBroadcast для работы с сообщениями в каждой комнате чата
