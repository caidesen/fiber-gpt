package main

import (
	"server/app/handler"
	"server/app/model"
	"server/pkg/config"
	"server/pkg/db"
	"server/pkg/server"
	"server/pkg/session"
	"server/pkg/storage"
)

func main() {
	config.SetupFormYaml()
	app := server.Setup()
	db.Setup()
	model.Migrate()
	storage.Setup()
	session.Setup(storage.GetStorage())
	apiGroup := app.Group("/api")
	handler.NewPageHandler(db.GetDB(), app)
	handler.NewAuthHandler(db.GetDB(), apiGroup.Group("/auth"))
	handler.NewSettingsHandler(db.GetDB(), apiGroup.Group("/settings"))
	handler.NewChatHandler(db.GetDB(), apiGroup.Group("/chat"))
	server.StartServer(app)
}
