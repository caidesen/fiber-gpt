package repo

import (
	"context"
	"server/app/model"
	"server/pkg/db"
)

func GetUserById(c context.Context, db db.Runner, userId uint) (*model.User, error) {
	var user model.User
	err := db.QueryRowContext(c, "SELECT id, nickname, username, password FROM users WHERE id = ?", userId).Scan(&user.Id, &user.Nickname, &user.Username, &user.Password)
	return &user, err
}
func GetUserByUsername(c context.Context, db db.Runner, username string) (*model.User, error) {
	var user model.User
	err := db.QueryRowContext(c, "SELECT id, nickname, username, password FROM users WHERE username = ?", username).Scan(&user.Id, &user.Nickname, &user.Username, &user.Password)
	return &user, err
}

func InsertUser(c context.Context, db db.Runner, user *model.User) error {
	_, err := db.ExecContext(c, "INSERT INTO users (nickname, username, password) VALUES (?, ?, ?)", user.Nickname, user.Username, user.Password)
	return err
}
func UpdateUser(c context.Context, db db.Runner, user *model.User) error {
	_, err := db.ExecContext(c, "UPDATE users SET nickname = ?, username = ?, password = ? WHERE id = ?", user.Nickname, user.Username, user.Password, user.Id)
	return err
}
func DeleteUser(c context.Context, db db.Runner, userId uint) error {
	_, err := db.ExecContext(c, "DELETE FROM users WHERE id = ?", userId)
	return err
}
