package handler

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"net/http"
	"server/pkg/session"
	"server/web"
)

type PageHandler struct {
	db *sql.DB
}

func NewPageHandler(db *sql.DB, a fiber.Router) *PageHandler {
	handler := PageHandler{
		db: db,
	}
	handler.setup(a)
	return &handler
}
func (h *PageHandler) renderPage(ctx *fiber.Ctx) error {
	file, _ := web.GetWebContent().ReadFile("dist/index.html")
	ctx.Response().Header.Set("Content-Type", "text/html; charset=utf-8")
	return ctx.Send(file)
}
func (h *PageHandler) setup(a fiber.Router) {
	a.Use("/assets", filesystem.New(filesystem.Config{
		Root:       http.FS(web.GetWebContent()),
		PathPrefix: "dist/assets",
	}))
	a.Use(favicon.New(favicon.Config{
		File:         "dist/vite.svg",
		FileSystem:   http.FS(web.GetWebContent()),
		URL:          "/vite.svg",
		CacheControl: "public, max-age=31536000",
	}))
	a.Get("/login", h.renderPage)
	a.Get("/register", func(ctx *fiber.Ctx) error {
		id, err := session.GetUserId(ctx)
		if err != nil || id == 0 {
			return ctx.Next()
		} else {
			return ctx.Redirect("/")
		}
	}, h.renderPage)
	a.Get("/*", func(ctx *fiber.Ctx) error {
		id, err := session.GetUserId(ctx)
		if err != nil || id == 0 {
			return ctx.Redirect("/login")
		}
		return ctx.Next()
	}, h.renderPage)
}
