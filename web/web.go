package web

import "embed"

//go:embed dist/*
var content embed.FS

func GetWebContent() embed.FS {
	return content
}
