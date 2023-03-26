package session

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store *session.Store

func Setup(storage fiber.Storage) {
	store = session.New(session.Config{
		Storage: storage,
	})
}
func GetSession(c *fiber.Ctx) (*session.Session, error) {
	return store.Get(c)
}
func GetUserId(c *fiber.Ctx) (uint, error) {
	sess, err := GetSession(c)
	if err != nil {
		return 0, err
	}
	uid := sess.Get("uid")
	if uid == nil {
		return 0, nil
	}
	return uid.(uint), nil
}
func AuthRequired(c *fiber.Ctx) error {
	id, err := GetUserId(c)
	if err != nil || id == 0 {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}
