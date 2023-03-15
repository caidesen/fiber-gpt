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
