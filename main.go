package telebot

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const (
	secThrottle = time.Second * 2
)

//ErrCongested means, that the buffer has no free space
var ErrCongested = errors.New("buffer congested")

//Bot is the config object
type Bot struct {
	accessToken string
	chatID      string
	msgBuf      chan string
	backlogSize int
	ErrHandler  func(error)
	LocalEcho   bool
}

//New creates a new telebot instance
func New(accessToken, chatID string, backlogSize int) *Bot {

	bot := &Bot{
		accessToken: accessToken,
		chatID:      chatID,
		msgBuf:      make(chan string, backlogSize),
		backlogSize: backlogSize,
		LocalEcho:   false,
		ErrHandler: func(e error) {
			// default error handler
			fmt.Println(fmt.Sprintf("%s\n", e.Error()))
		},
	}

	go worker(bot)
	return bot
}

//SetErrHandler sets a custom error handler
func (bot *Bot) SetErrHandler(f func(error)) {
	bot.ErrHandler = f
}

//Bye frees ressources
func (bot *Bot) Bye() {
	close(bot.msgBuf)
}

//Send sends a message to the configured bot
func (bot *Bot) Send(message string) error {

	// perform a non blocking channel write
	select {
	case bot.msgBuf <- message:
		return nil
	default:
		// discard message
		return ErrCongested
	}
}

func worker(bot *Bot) {

	for message := range bot.msgBuf {
		sendMessage(bot, message)
	}
}

func sendMessage(tb *Bot, message string) {

	if tb.LocalEcho {
		oneSecond := float64(time.Second)
		// wait for a random of 0 and 1.5s
		// (never causing timeout errors)
		secs := oneSecond * 1.5 * rand.Float64()
		time.Sleep(time.Duration(secs))
		fmt.Printf("telebot: %s\n", message)
		return
	}

	start := time.Now()
	deadline := start.Add(secThrottle)

	err := sendToTelegramBot(tb, message)

	if (err != nil) && (tb.ErrHandler != nil) {
		tb.ErrHandler(err)
		return
	}

	elapsed := time.Since(start)

	if elapsed.Seconds() < float64(secThrottle) {
		time.Sleep(time.Until(deadline))
	}
}
