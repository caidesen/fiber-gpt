package session

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"server/pkg/storage"
)

var store *session.Store

func Setup() {
	store = session.New(session.Config{
		Storage: storage.New(storage.Config{
			Database: "session.db",
		}),
	})
}
func GetSession(c *fiber.Ctx) (*session.Session, error) {
	return store.Get(c)
}
