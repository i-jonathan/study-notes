package main

import (
	"fmt"
	goTel "github.com/yoruba-codigy/goTelegram"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

var bot goTel.Bot
var notesList map[int]*pendingNotes
var pubCategories []string
var db *gorm.DB

func main() {
	db = initDatabase()
	notesList = make(map[int]*pendingNotes)
	pubCategories = []string{"Article", "Bible", "Broadcast", "Brochure", "Meetings & Conventions", "Magazines",
		"Special Programs"}

	var err error
	bot, err = goTel.NewBot(os.Getenv("bot_token"))
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(bot.Me.Firstname, bot.Me.ID)
	bot.SetHandler(handler)

	log.Println("Starting Server...")
	err = http.ListenAndServe(":" + os.Getenv("port"), http.HandlerFunc(bot.UpdateHandler))

	if err != nil {
		log.Println(err)
		log.Fatalln(err)
	}
}

