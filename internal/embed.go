package internal

import "embed"

//go:embed pb/*/*/*.swagger.json
var SwaggerContent embed.FS

//go:embed ui/dist/*
var FrontendContent embed.FS
