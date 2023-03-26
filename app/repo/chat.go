package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"server/app/model"
	"server/pkg/db"
	"time"
)

func InsertChat(c context.Context, db db.Runner, chat *model.Chat) error {
	settings, _ := json.Marshal(chat.Settings)
	now := time.Now()
	chat.CreatedAt = now
	chat.UpdatedAt = now
	sql := "INSERT INTO chats (user_id, settings, created_at, updated_at) VALUES ($1, $2, $3, $4) returning id"
	err := db.QueryRowContext(c, sql, chat.UserId, settings, now, now).Scan(&chat.Id)
	return err
}
func FindChatsWithLatestQuestionByUserId(c context.Context, db db.Runner, userId uint) ([]*model.Chat, error) {
	var chats []*model.Chat
	rows, err := db.QueryContext(c, `
		SELECT c.id, c.user_id, c.settings, c.created_at, c.updated_at, q.id, q.q, q.a, q.status, q.created_at, q.updated_at
		FROM chats c
		LEFT JOIN questions q ON q.id = (
		    	SELECT id FROM questions WHERE chat_id = c.id ORDER BY created_at DESC LIMIT 1
		)
		WHERE c.user_id = $1
		ORDER BY q.created_at DESC
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var chat model.Chat
		var qId sql.NullInt64
		var q sql.NullString
		var a sql.NullString
		var createdAt sql.NullTime
		var updatedAt sql.NullTime
		var status sql.NullString
		var settings []byte
		if err := rows.Scan(&chat.Id, &chat.UserId, &settings, &chat.CreatedAt, &chat.UpdatedAt, &qId, &q, &a, &status, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal(settings, &chat.Settings)
		chat.Questions = make([]*model.Question, 0)
		if qId.Valid {
			var question model.Question
			question.Id = uint(qId.Int64)
			question.Q = q.String
			question.A = a.String
			question.Status = status.String
			question.CreatedAt = createdAt.Time
			question.UpdatedAt = updatedAt.Time
			chat.Questions = append(chat.Questions, &question)
		}
		chats = append(chats, &chat)
	}
	return chats, err
}
func FindOneChatById(c context.Context, db db.Runner, chatId uint) (*model.Chat, error) {
	var chat model.Chat
	var settings []byte
	err := db.
		QueryRowContext(c, "SELECT id, user_id, settings, created_at, updated_at FROM chats WHERE id = $1", chatId).
		Scan(&chat.Id, &chat.UserId, &settings, &chat.CreatedAt, &chat.UpdatedAt)
	json.Unmarshal(settings, &chat.Settings)
	return &chat, err
}
func FindQuestionsByChatId(c context.Context, db db.Runner, chatId uint) ([]*model.Question, error) {
	var questions []*model.Question
	rows, err := db.QueryContext(c, "SELECT id, chat_id, q, a, status, created_at, updated_at FROM questions WHERE chat_id = $1", chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return questions, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var question model.Question
		if err := rows.Scan(&question.Id, &question.ChatId, &question.Q, &question.A, &question.Status, &question.CreatedAt, &question.UpdatedAt); err != nil {
			return nil, err
		}
		questions = append(questions, &question)
	}
	return questions, err
}
func FindOneQuestionById(c context.Context, db db.Runner, id uint) (*model.Question, error) {
	var question model.Question
	err := db.
		QueryRowContext(c, "SELECT id, chat_id, q, a, created_at, updated_at FROM questions WHERE id = $1", id).
		Scan(&question.Id, &question.ChatId, &question.Q, &question.A, &question.CreatedAt, &question.UpdatedAt)
	return &question, err
}
func InsertQuestion(c context.Context, db db.Runner, question *model.Question) error {
	now := time.Now()
	question.CreatedAt = now
	question.UpdatedAt = now
	return db.
		QueryRowContext(
			c,
			"INSERT INTO questions (chat_id, q, a,status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
			question.ChatId, question.Q, question.A, question.Status, now, now).
		Scan(&question.Id)
}
func UpdateQuestion(c context.Context, db db.Runner, question *model.Question) error {
	now := time.Now()
	_, err := db.ExecContext(
		c,
		"UPDATE questions SET a = $1, updated_at = $2, status = $3 WHERE id = $4",
		question.A, now, question.Status, question.Id,
	)
	return err
}
