package main

import (
	"fmt"
	goTel "github.com/yoruba-codigy/goTelegram"
	"log"
	"net/http"
	"os"
	"time"
)

var bot goTel.Bot
var notesList map[int]*pendingNotes

func main() {
	notesList = make(map[int]*pendingNotes)

	var err error
	bot, err = goTel.NewBot(os.Getenv("bot_token"))
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(bot.Me.Firstname, bot.Me.ID)
	bot.SetHandler(handler)

	log.Println(time.Now(), ": Starting Server...")
	err = http.ListenAndServe(":" + os.Getenv("port"), http.HandlerFunc(bot.UpdateHandler))

	if err != nil {
		log.Println(time.Now(), ":", err)
		log.Fatalln(err)
	}
}

