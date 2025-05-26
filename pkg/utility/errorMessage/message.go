package errorMessage

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	TH = "TH"
	EN = "EN"
)

type ErrorMessage struct {
	Code    int     `json:"code"`
	Message Message `json:"message"` // index to lang
}

type Message map[string]string

func (m Message) Language(l string) string {
	var message string
	for k, v := range m {
		if k == l {
			return v
		}
		if strings.ToUpper(k) == EN {
			message = v
		}
		if message == "" {
			message = v
		}
	}
	return message
}

func (m Message) LanguageByContext(ctx context.Context) string {
	language := EN
	l := strings.ToUpper(fmt.Sprintf("%v", ctx.Value(fiber.HeaderAcceptLanguage)))
	if l != "" {
		language = l
	}
	var message string
	for k, v := range m {
		if k == language {
			return v
		}
		if strings.ToUpper(k) == EN {
			message = v
		}
		if message == "" {
			message = v
		}
	}
	return message
}
