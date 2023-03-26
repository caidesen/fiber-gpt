package storage

import (
	"github.com/gofiber/storage/memory"
	"time"
)

var storage *memory.Storage

func GetStorage() *memory.Storage {
	return storage
}

func Setup() {
	storage = memory.New(memory.Config{
		GCInterval: 1 * time.Minute,
	})
}
