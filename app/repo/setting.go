package repo

import (
	"context"
	"server/app/model"
	"server/pkg/db"
)

func FindSettings(c context.Context, db db.Runner, userId uint) ([]*model.Setting, error) {
	var settings []*model.Setting
	rows, err := db.QueryContext(c, "SELECT key, value FROM settings WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var setting model.Setting
		setting.UserId = userId
		if err := rows.Scan(&setting.Key, &setting.Value); err != nil {
			return nil, err
		}
		settings = append(settings, &setting)
	}
	return settings, err
}

func UpsertSettings(c context.Context, db db.Runner, settings []*model.Setting) error {
	for _, setting := range settings {
		_, err := db.ExecContext(
			c,
			"INSERT OR REPLACE INTO settings (user_id, key, value) VALUES ($1, $2, $3)",
			setting.UserId, setting.Key, setting.Value,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
