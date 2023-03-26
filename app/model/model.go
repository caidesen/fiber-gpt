package model

import (
	"golang.org/x/crypto/bcrypt"
	"server/pkg/db"
	"time"
)

type User struct {
	Id       uint
	Nickname string
	Username string
	Password string
}
type Setting struct {
	UserId uint
	Key    string
	Value  string
}

type Chat struct {
	Id        uint
	UserId    uint
	Questions []*Question
	Settings  struct {
		MaxToken    int
		Temperature float32
	}
	CreatedAt time.Time
	UpdatedAt time.Time
}
type Question struct {
	Id        uint
	ChatId    uint
	Q         string
	A         string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var createTable = `
CREATE TABLE IF NOT EXISTS users
(
    id INTEGER PRIMARY KEY, 
	nickname TEXT, 
	username TEXT, 
	password TEXT
);
CREATE TABLE IF NOT EXISTS settings
(
    user_id INTEGER,
    key TEXT,
    value TEXT,
    primary key (user_id, key)
);
CREATE TABLE IF NOT EXISTS chats 
(
   	id INTEGER PRIMARY KEY,
   	user_id INTEGER,
   	settings TEXT,
   	created_at DATE,
    updated_at DATE 
);
CREATE TABLE IF NOT EXISTS questions (
   id INTEGER PRIMARY KEY,
   chat_id INTEGER,
   q TEXT,
   a TEXT,
   status TEXT,
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
		_, err = db.Exec("INSERT INTO users (nickname, username, password) VALUES ($1, $2, $3)", "admin", "admin", password)
		if err != nil {
			panic(err)
		}
	}
}
