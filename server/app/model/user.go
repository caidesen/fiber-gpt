package model

import (
	"golang.org/x/crypto/bcrypt"
	"server/pkg/db"
	"time"
)

type User struct {
	Id       uint   `db:"id"`
	Nickname string `db:"nickname"`
	Username string `db:"username"`
	Password string `db:"password"`
}
type GptSetting struct {
	UserId      uint    `db:"user_id"`
	ApiKey      string  `db:"api_key"`
	Temperature float32 `db:"temperature"`
	MaxToken    int     `db:"max_token"`
}
type Chat struct {
	Id          uint    `db:"id"`
	UserId      uint    `db:"user_id"`
	MaxToken    int     `db:"max_token"`
	Temperature float32 `db:"Temperature"`
}
type Question struct {
	Id        uint      `db:"id"`
	Q         string    `db:"q"`
	A         string    `db:"a"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

var createTable = `
CREATE TABLE IF NOT EXISTS users
(
    id INTEGER PRIMARY KEY, 
	nickname TEXT, 
	username TEXT, 
	password TEXT
);
CREATE TABLE IF NOT EXISTS gpt_settings
(
   user_id INTEGER PRIMARY KEY,
   api_key TEXT,
   temperature TEXT,
   max_token INTEGER 
);
CREATE TABLE IF NOT EXISTS chats 
(
   	id INTEGER PRIMARY KEY,
   	user_id INTEGER,
   	max_token INTEGER,
   	temperature TEXT
);
CREATE TABLE IF NOT EXISTS questions (
   id INTEGER PRIMARY KEY,
   q TEXT,
   a TEXT,
   created_at DATE,
   updated_at DATE 
)	
`

func Migrate() {
	db := db.GetDB()
	_, err := db.Exec(createTable)
	// query user count
	if err != nil {
		panic(err)
	}
	var count int
	if err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		panic(err)
	}
	if count == 0 {
		password, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.MinCost)
		if err != nil {
			return
		}
		_, err = db.Exec("INSERT INTO users (nickname, username, password) VALUES (?, ?, ?)", "admin", "admin", password)
		if err != nil {
			panic(err)
		}
	}

}
