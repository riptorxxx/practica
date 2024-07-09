package main

import (
	"database/sql"
	. "fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

// Реализуем подключение к БД
func initDatabase() *sql.DB {
	psqlInfo := Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.PgHost, cfg.PgPort, cfg.PgUser, cfg.PgPass, cfg.PgBase)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	return db
}

// Создаём все необходимое в postgres для работы приложения.
func createSqlObj() {
	var err error
	createUsersTableQuery := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        login VARCHAR(50) UNIQUE NOT NULL,
        email VARCHAR(50) UNIQUE NOT NULL,
        phone VARCHAR(15) UNIQUE NOT NULL,
        password TEXT NOT NULL,
		token TEXT,
        uid UUID DEFAULT gen_random_uuid() UNIQUE,
        create_date TIMESTAMP DEFAULT current_timestamp
    )`

	_, err = db.Exec(createUsersTableQuery)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	createChatsTableQuery := `
    CREATE TABLE IF NOT EXISTS chats (
		id SERIAL PRIMARY KEY,
        uid UUID DEFAULT gen_random_uuid(),
        name VARCHAR(64) NOT NULL,
        lcc TIMESTAMP NOT NULL,
        cypher VARCHAR(64) NOT NULL,
        user_uid UUID REFERENCES users(uid)
    )`

	_, err = db.Exec(createChatsTableQuery)
	if err != nil {
		log.Fatal("Failed to create chats table:", err)
	}
}

func deleteExpiredChats() {
	for {
		// Сначала выбираем устаревшие чаты
		rows, err := db.Query("SELECT name, uid FROM chats WHERE lcc < NOW()")
		if err != nil {
			log.Printf("Error selecting expired chats: %v", err)
			// time.Sleep(1 * time.Minute)
			continue
		}

		// Записываем информацию о чатах в лог и собираем UUID для удаления
		var uuids []string
		for rows.Next() {
			var name string
			var uuid string
			if err := rows.Scan(&name, &uuid); err != nil {
				log.Printf("Error scanning chat row: %v", err)
				continue
			}
			log.Printf("Deleting chat: Name=%s, UUID=%s", name, uuid)
			uuids = append(uuids, uuid)
		}
		rows.Close()

		// Удаляем устаревшие чаты
		if len(uuids) > 0 {
			result, err := db.Exec("DELETE FROM chats WHERE lcc < NOW()")
			if err != nil {
				log.Printf("Error deleting expired chats: %v", err)
			} else {
				rowsAffected, _ := result.RowsAffected()
				log.Printf("Expired chats deleted: %d", rowsAffected)
			}
		}
		// time.Sleep(1 * time.Minute)
	}
}

// В PostgreSQL можно создать фоновую задачу, которая будет периодически проверять таблицу чатов
// и удалять устаревшие записи. Один из способов реализации — использование cron-задач,
// например, с помощью расширения pg_cron.

// #### Шаги для реализации:
// 1. **Установите расширение pg_cron:**
//     CREATE EXTENSION pg_cron;

// Создайте функцию для удаления устаревших чатов:в:**
//     CREATE OR REPLACE FUNCTION delete_expired_chats() RETURNS void AS $$
//     BEGIN
//         DELETE FROM chats WHERE lcc < NOW();
//     END;
//     $$ LANGUAGE plpgsql;

// 3. **Создайте cron-задачу для периодического выполнения функции:**
//     SELECT cron.schedule('delete_expired_chats', '0 * * * *', 'SELECT delete_expired_chats();');
//     Эта задача будет выполняться каждый час и удалять устаревшие чаты.
