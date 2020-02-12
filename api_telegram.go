package telebot

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL = "https://api.telegram.org/bot"
)

func sendToTelegramBot(bot *Bot, message string) error {

	rawURL := fmt.Sprintf("%s%s/sendMessage", baseURL, bot.accessToken)
	url, _ := url.Parse(rawURL)
	q := url.Query()
	q.Set("chat_id", bot.chatID)
	q.Set("text", message)
	url.RawQuery = q.Encode()

	client := http.Client{Timeout: 5 * time.Second}

	ddd := url.String()
	fmt.Println(ddd)
	req, err := http.NewRequest("GET", ddd, nil)

	if err != nil {
		return err
	}

	res, err := client.Do(req)

	if err != nil {
		return err
	}

	sendMessageJSON := struct {
		OK bool `json:"ok"`
	}{}

	err = json.NewDecoder(res.Body).Decode(&sendMessageJSON)

	if err != nil {
		return err
	}

	if !sendMessageJSON.OK {
		return errors.New("failed to send message")
	}

	return nil
}
