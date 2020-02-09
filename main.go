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

// Telebot is the config object
type Telebot struct {
	accessToken string
	chatID      string
	msgBuf      chan string
	backlogSize int
	ErrHandler  func(error)
	LocalEcho   bool
}

//New creates a new telebot instance
func New(accessToken, chatID string, backlogSize int) *Telebot {

	telebot := &Telebot{
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

	go worker(telebot)
	return telebot
}

//SetErrHandler sets a custom error handler
func (tb *Telebot) SetErrHandler(f func(error)) {
	tb.ErrHandler = f
}

//Bye frees ressources
func (tb *Telebot) Bye() {
	close(tb.msgBuf)
}

//Send sends a message to the configured bot
func (tb *Telebot) Send(message string) error {

	// perform a non blocking channel write
	select {
	case tb.msgBuf <- message:
		return nil
	default:
		// discard message
		return ErrCongested
	}
}

func worker(tb *Telebot) {

	for message := range tb.msgBuf {
		sendMessage(tb, message)
	}
}

func sendMessage(tb *Telebot, message string) {

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
