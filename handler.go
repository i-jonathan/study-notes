package main

import (
	goTel "github.com/yoruba-codigy/goTelegram"
	"log"
	"math/rand"
	"time"
)

func handler(update goTel.Update) {
	bot.DeleteKeyboard()
	switch update.Type {
	case "text":
		// process text updates
		processText(update)
	case "callback":
		//process callback
		processCallBack(update)
		return
	}
}

func processText(update goTel.Update) {
	// update.Text starts with #, process tags
	log.Println(update.Command)
	switch update.Command {
	case "/start":
		mainMenu(update)
	}
}

func processCallBack(update goTel.Update) {
	switch update.CallbackQuery.Data {
	case "addNote":
		newNote := studyNotes{
			CreatedAt: time.Now(),
			UserId:      update.CallbackQuery.From.ID,
		}

		message, err := bot.EditMessage(update.CallbackQuery.Message, "Give your Note an expressive title:")

		if err != nil {
			log.Println("Error when sending stage 0 message", err)
			return
		}

		processNote := pendingNotes{
			Stages:       5,
			CurrentStage: 0,
			Message:      message,
			Data:         newNote,
		}

		notesList[update.CallbackQuery.From.ID] = &processNote
	}
}

func mainMenu(update goTel.Update) {
	greetings := []string{"Bonjour", "Hola", "Konnichiwa", "Hey", "Hello"}
	rand.Seed(time.Now().Unix())
	text := greetings[rand.Intn(len(greetings))] + ", "
	bot.AddButton("Create Note", "addNote")
	bot.AddButton("List Notes", "listNotes")
	bot.AddButton("List Tags", "listTags")
	bot.AddButton("List By Tags", "tagList")
	bot.AddButton("List By Category", "categoryList")
	bot.MakeKeyboard(1)

	if update.Type == "text" {
		if update.Message.Chat.Type != "private" {
			_, err := bot.SendMessage("Can't use this in a group.", update.Message.Chat)
			if err != nil {
				log.Println("error in sending message in menu", err)
				return
			}
		}

		if update.Message.From.ID == 726094693 {
			moniker := []string{"Babe", "Baby", "Asunke mi", "🌝", "Yoyo", "Clown"}
			rand.Seed(time.Now().Unix())
			text += moniker[rand.Intn(len(moniker))]
		} else {
			text += update.Message.From.Firstname
		}

		_, err := bot.SendMessage(text, update.Message.Chat)
		if err != nil {
			log.Println("error when sending welcome message", err)
		}
	} else {
		if update.CallbackQuery.Message.Chat.Type != "private" {
			_, err := bot.SendMessage("Can't use this in a group.", update.Message.Chat)
			if err != nil {
				log.Println("error in sending message in menu", err)
				return
			}
		}

		if update.Message.From.ID == 726094693 {
			moniker := []string{"Babe", "Baby", "Asunke mi", "🌝", "Yoyo", "Clown"}
			rand.Seed(time.Now().Unix())
			text += moniker[rand.Intn(len(moniker))]
		} else {
			text += update.CallbackQuery.From.Firstname
		}
		_, err := bot.SendMessage(text, update.Message.Chat)
		if err != nil {
			log.Println("error when sending welcome message", err)
		}
	}
}