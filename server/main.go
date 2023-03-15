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
	config.SetupFormYaml("config.yaml")
	app := server.Setup()
	db.Setup()
	model.Migrate()
	storage.Setup()
	session.Setup(storage.GetStorage())
	apiGroup := app.Group("/api")
	handler.NewAuthHandler(db.GetDB(), apiGroup)
	server.StartServer(app)
}
