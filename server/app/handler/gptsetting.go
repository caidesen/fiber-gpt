package handler

import (
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"server/app/model"
	"server/app/repo"
	"server/pkg/server"
	"server/pkg/session"
)

type GptSettingHandler struct {
	db *sql.DB
}

func NewGptSettingHandler(db *sql.DB, a fiber.Router) *GptSettingHandler {
	handler := GptSettingHandler{db: db}
	handler.setup(a)
	return &handler
}

func (h *GptSettingHandler) setup(a fiber.Router) {
	g := a.Group("/gptsetting")
	g.Post("update", server.WithJSONBody(h.update))
	g.Post("get", h.get)
}

type gptSettingVo struct {
	ApiKey      string  `json:"apiKey"`
	Temperature float32 `json:"temperature"`
	MaxToken    int     `json:"maxToken"`
}
type updateGptSettingInput struct {
	ApiKey      *string  `json:"apiKey"`
	Temperature *float32 `json:"temperature"`
	MaxToken    *int     `json:"maxToken"`
}

func (h *GptSettingHandler) get(c *fiber.Ctx) error {
	userId, err := session.GetUserId(c)
	if err != nil {
		return err
	}
	setting, err := repo.GetGptSettingByUserId(c.Context(), h.db, userId)
	if errors.Is(err, sql.ErrNoRows) {
		setting = &model.GptSetting{
			UserId:      userId,
			ApiKey:      "",
			Temperature: 0.7,
			MaxToken:    100,
		}
		_, err = repo.UpsertGptSetting(c.Context(), h.db, setting)
	}
	if err != nil {
		return err
	}
	vo := gptSettingVo{
		ApiKey:      setting.ApiKey,
		Temperature: setting.Temperature,
		MaxToken:    setting.MaxToken,
	}
	return c.JSON(vo)
}

// updateGptSetting updates the gpt setting of the current user
func (h *GptSettingHandler) update(c *fiber.Ctx, input *updateGptSettingInput) error {
	userId, err := session.GetUserId(c)
	if err != nil {
		return err
	}
	setting, err := repo.GetGptSettingByUserId(c.Context(), h.db, userId)
	if err != nil {
		return err
	}
	if input.ApiKey != nil {
		setting.ApiKey = *input.ApiKey
	}
	if input.Temperature != nil {
		setting.Temperature = *input.Temperature
	}
	if input.MaxToken != nil {
		setting.MaxToken = *input.MaxToken
	}
	_, err = repo.UpsertGptSetting(c.Context(), h.db, setting)
	return err
}
