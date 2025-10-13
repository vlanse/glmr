package main

import (
	"context"
	"fmt"
	"time"
)

func init() {
	time.Local = time.UTC
}

func main() {
	app := NewApp()
	if err := app.Run(context.Background()); err != nil {
		fmt.Println(err.Error())
	}
}
