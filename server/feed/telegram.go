package feed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var ErrTelegramBotTokenNotSet = errors.New("TELEGRAM_BOT_TOKEN environment variable not set")
var ErrTelegramBotAPIStatusCode = errors.New("telegram API error")
var ErrTelegramNotOKResponse = errors.New("telegram API returned a not OK status")

type (
	Chat struct {
		ID    int64  `json:"id"`
		Title string `json:"title,omitempty"`
		Type  string `json:"type"`
	}
	Message struct {
		MessageID int    `json:"message_id"`
		Date      int64  `json:"date"`
		Text      string `json:"text,omitempty"`
		Chat      *Chat  `json:"chat,omitempty"`
	}
	Update struct {
		UpdateID    int      `json:"update_id"`
		Message     *Message `json:"message,omitempty"`
		ChannelPost *Message `json:"channel_post,omitempty"`
	}
	TelegramResponse struct {
		OK      bool     `json:"ok"`
		Results []Update `json:"result"`
	}
)

func buildTelegramURL(chatID string, limit uint) (string, error) {
	token, found := os.LookupEnv("TELEGRAM_BOT_TOKEN")
	if !found || token == "" {
		return "", ErrTelegramBotTokenNotSet
	}

	u := &url.URL{
		Scheme: "https",
		Host:   "api.telegram.org",
		Path:   fmt.Sprintf("/bot%s/getUpdates", token),
	}

	params := u.Query()
	params.Set("chat_id", chatID)
	params.Set("limit", strconv.Itoa(int(limit)))
	u.RawQuery = params.Encode()

	return u.String(), nil

}

func fetchTelegramChannelMessages(
	ctx context.Context,
	chatID string,
	limit uint,
) ([]string, error) {
	telegramURL, err := buildTelegramURL(chatID, limit)
	if err != nil {
		return nil, err
	}

	reqCtx, reqCancel := context.WithTimeout(ctx, 10*time.Second)
	defer reqCancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, telegramURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("could not close the request body successfully", "errro", err.Error())
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Join(
			fmt.Errorf("telegram API error: status %d", resp.StatusCode),
			ErrTelegramBotAPIStatusCode,
		)
	}

	var response TelegramResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if !response.OK {
		return nil, ErrTelegramNotOKResponse
	}

	var messages []string
	for _, result := range response.Results {
		text := ""
		if result.Message != nil && result.Message.Text != "" {
			text = result.Message.Text
		} else if result.ChannelPost != nil && result.ChannelPost.Text != "" {
			text = result.ChannelPost.Text
		}

		if text != "" {
			messages = append(messages, text)
		}
	}

	return messages, nil
}
