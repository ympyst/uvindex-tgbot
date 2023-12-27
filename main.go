package main

import (
	"context"
	"github.com/ympyst/uvindex-tgbot/telegram"
)

func main() {
	t := telegram.NewTelegram()
	ctx := context.Background()

	t.Start(ctx)
}

