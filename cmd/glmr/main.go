package main

import (
	"context"
	"log"
	"time"
)

func init() {
	time.Local = time.UTC
}

func main() {
	app := NewApp()
	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
