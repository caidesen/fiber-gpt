package repo

import (
	"context"
	"server/app/model"
	"server/pkg/db"
)

func FindOneUserById(c context.Context, db db.Runner, userId uint) (*model.User, error) {
	var user model.User
	err := db.QueryRowContext(c, "SELECT id, nickname, username, password FROM users WHERE id = $1", userId).Scan(&user.Id, &user.Nickname, &user.Username, &user.Password)
	return &user, err
}
func FindOneUserByUsername(c context.Context, db db.Runner, username string) (*model.User, error) {
	var user model.User
	err := db.QueryRowContext(c, "SELECT id, nickname, username, password FROM users WHERE username = $1", username).Scan(&user.Id, &user.Nickname, &user.Username, &user.Password)
	return &user, err
}

func InsertUser(c context.Context, db db.Runner, user *model.User) error {
	_, err := db.ExecContext(c, "INSERT INTO users (nickname, username, password) VALUES ($1, $2, $3)", user.Nickname, user.Username, user.Password)
	return err
}
func UpdateUser(c context.Context, db db.Runner, user *model.User) error {
	_, err := db.ExecContext(c, "UPDATE users SET nickname = $1, username = $2, password = $3 WHERE id = $4", user.Nickname, user.Username, user.Password, user.Id)
	return err
}
func DeleteUser(c context.Context, db db.Runner, userId uint) error {
	_, err := db.ExecContext(c, "DELETE FROM users WHERE id = $1", userId)
	return err
}
