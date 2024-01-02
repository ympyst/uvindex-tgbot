package main

import (
	"context"
	"github.com/ympyst/uvindex-tgbot/app"
)

func main() {
	a := app.NewApp()
	ctx := context.Background()

	a.Start(ctx)
}

