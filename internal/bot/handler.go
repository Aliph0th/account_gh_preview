package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"tokens-bot/internal/database"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const SUCCESS = "<tg-emoji emoji-id=\"5435881193033711915\">😊</tg-emoji>"
const ERROR = "<tg-emoji emoji-id=\"5462963147973340850\">😞</tg-emoji>"

type Handler struct {
	db *database.DB
}

func New(db *database.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) RegisterHandlers(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/tokens", bot.MatchTypePrefix, h.handleTokensCommand)
}

func (h *Handler) handleTokensCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	parts := strings.Fields(update.Message.Text)
	if len(parts) != 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      fmt.Sprintf("%v Incorrect format. Use: /tokens <uint64>", ERROR),
			ParseMode: "HTML",
		})
		return
	}

	_, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      fmt.Sprintf("%v Invalid number. Please provide a valid number value.", ERROR),
			ParseMode: "HTML",
		})
		return
	}

	if err := h.db.SetTokens(parts[1]); err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      fmt.Sprintf("%v Database error: %v", ERROR, err),
			ParseMode: "HTML",
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      fmt.Sprintf("%v Tokens updated: %s", SUCCESS, parts[1]),
		ParseMode: "HTML",
	})
}
