package main

import (
	"fmt"
	goTel "github.com/yoruba-codigy/goTelegram"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

var (
	bot goTel.Bot
	notesList map[int]*pendingNotes
	pubCategories []string
	db *gorm.DB
	err error
)

func main() {
	db = initDatabase()
	notesList = make(map[int]*pendingNotes)
	pubCategories = []string{"Article", "Bible", "Broadcast", "Brochure", "Meetings & Conventions", "Magazines",
		"Special Programs"}

	bot, err = goTel.NewBot(os.Getenv("bot_token"))
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(bot.Me.Firstname, bot.Me.ID)
	bot.SetHandler(handler)

	log.Println("Starting Server...")
	err = http.ListenAndServe(":"+os.Getenv("PORT"), http.HandlerFunc(bot.UpdateHandler))

	if err != nil {
		log.Println(err)
		log.Fatalln(err)
	}
}
