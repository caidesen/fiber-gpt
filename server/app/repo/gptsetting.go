package repo

import (
	"context"
	"server/app/model"
	"server/pkg/db"
)

func GetGptSettingByUserId(c context.Context, db db.Runner, userId uint) (*model.GptSetting, error) {
	var gptSetting model.GptSetting
	err := db.
		QueryRowContext(c, "SELECT user_id, api_key, temperature, max_token FROM gpt_settings WHERE user_id = ?", userId).
		Scan(&gptSetting.UserId, &gptSetting.ApiKey, &gptSetting.Temperature, &gptSetting.MaxToken)
	return &gptSetting, err
}
func UpsertGptSetting(c context.Context, db db.Runner, setting *model.GptSetting) (*model.GptSetting, error) {
	_, err := db.ExecContext(
		c,
		"INSERT OR REPLACE INTO gpt_settings (user_id, api_key, temperature, max_token) VALUES (?, ?, ?, ?)",
		setting.UserId, setting.ApiKey, setting.Temperature, setting.MaxToken,
	)
	if err != nil {
		return nil, err
	}
	return setting, nil
}
