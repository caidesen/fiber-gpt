package handler

import (
	"context"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"server/app/model"
	"server/app/repo"
	"server/pkg/db"
	"server/pkg/server"
	"server/pkg/session"
	"strconv"
)

type SettingsHandler struct {
	db *sql.DB
}

func NewSettingsHandler(db *sql.DB, a fiber.Router) *SettingsHandler {
	handler := SettingsHandler{db: db}
	handler.setup(a)
	return &handler
}
func (h *SettingsHandler) setup(a fiber.Router) {
	a.Use(session.AuthRequired)
	a.Post("update", server.WithJSONBody(h.update))
	a.Post("get", h.get)
}
func (h *SettingsHandler) get(c *fiber.Ctx) error {
	userId, err := session.GetUserId(c)
	if err != nil {
		return err
	}
	settings, err := GetUserSettings(c.Context(), h.db, userId)
	if err != nil {
		return err
	}
	return c.JSON(settings)
}
func (h *SettingsHandler) update(c *fiber.Ctx, input *UserSettings) error {
	userId, err := session.GetUserId(c)
	if err != nil {
		return err
	}
	if err = SetUserSettings(c.Context(), h.db, userId, input); err != nil {
		return err
	}
	return c.JSON(&fiber.Map{
		"success": true,
	})
}

type UserSettings struct {
	ApiKey      string  `json:"apiKey"`
	MaxToken    int     `json:"maxToken"`
	Temperature float32 `json:"temperature"`
}

func GetUserSettings(c context.Context, db db.Runner, userId uint) (*UserSettings, error) {
	settings, err := repo.FindSettings(c, db, userId)
	if err != nil {
		return nil, err
	}
	userSettings := UserSettings{
		ApiKey:      "",
		MaxToken:    512,
		Temperature: 0.7,
	}
	for _, setting := range settings {
		switch setting.Key {
		case "apiKey":
			userSettings.ApiKey = setting.Value
		case "maxToken":
			mt, _ := strconv.Atoi(setting.Value)
			if mt != 0 {
				userSettings.MaxToken = mt
			}
		case "temperature":
			t, _ := strconv.ParseFloat(setting.Value, 32)
			if t != 0 {
				userSettings.Temperature = float32(t)
			}
		}
	}
	return &userSettings, nil
}
func SetUserSettings(c context.Context, db db.Runner, userId uint, UserSettings *UserSettings) error {
	settings := make([]*model.Setting, 0)
	if UserSettings.ApiKey != "" {
		settings = append(settings, &model.Setting{UserId: userId, Key: "apiKey", Value: UserSettings.ApiKey})
	}
	if UserSettings.MaxToken != 0 {
		settings = append(settings, &model.Setting{UserId: userId, Key: "maxToken", Value: strconv.Itoa(UserSettings.MaxToken)})
	}
	if UserSettings.Temperature != 0 {
		settings = append(settings, &model.Setting{UserId: userId, Key: "temperature", Value: strconv.FormatFloat(float64(UserSettings.Temperature), 'f', 2, 32)})
	}
	return repo.UpsertSettings(c, db, settings)
}
