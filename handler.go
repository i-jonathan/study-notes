package main

import (
	"fmt"
	goTel "github.com/yoruba-codigy/goTelegram"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func handler(update goTel.Update) {
	bot.DeleteKeyboard()
	switch update.Type {
	case "text":
		// process text updates
		if len(update.Command) == 0 {
			// process raw text
			processRawText(update)
			return
		}
		processCommand(update)
	case "callback":
		//process callback
		processCallBack(update)
	}
}

func processRawText(update goTel.Update) {
	//check for text that starts with #
	//check for notes being added
	currentUserNote := notesList[update.Message.From.ID]
	if currentUserNote != nil {
		handleNoteQuestions(update, currentUserNote)
		return
	}
}

func processCommand(update goTel.Update) {
	// update.Text starts with #, process tags
	log.Println(update.Command)
	switch update.Command {
	case "/start":
		mainMenu(update)
	}
}

func processCallBack(update goTel.Update) {
	callBack := update.CallbackQuery.Data
	if strings.HasPrefix(callBack, "listNotes") {
		callBack = "listNotes"
	} else if strings.HasPrefix(callBack, "note") {
		callBack = "note"
	} else if strings.HasPrefix(callBack, "delete") {
		callBack = "delete"
	} else if strings.HasPrefix(callBack, "deleteConfirm") {
		callBack = "deleteConfirm"
	}
	switch callBack {
	case "addNote":
		newNote := studyNote{
			CreatedAt: time.Now(),
			UserId:    update.CallbackQuery.From.ID,
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
	case "mainMenu":
		mainMenu(update)
	case "addNoteOk":
		// run function to insert note in DB
		created := createNote(notesList[update.CallbackQuery.From.ID].Data)
		text := "An error occurred."
		if created {
			text = "Your Note has been Created."
		}
		bot.AddButton("Menu", "mainMenu")
		bot.MakeKeyboard(1)
		_, err := bot.EditMessage(update.CallbackQuery.Message, text)
		if err != nil {
			log.Println(err)
			mainMenu(update)
		}
		delete(notesList, update.CallbackQuery.From.ID)
	case "bail":
		mainMenu(update)
		currentNote := notesList[update.CallbackQuery.From.ID]
		if currentNote != nil {
			delete(notesList, update.CallbackQuery.From.ID)
		}
	case "listNotes":
		text := listAllNotes(update.CallbackQuery.Data)
		_, err := bot.EditMessage(update.CallbackQuery.Message, text)
		if err != nil {
			log.Println(err)
		}
	case "note":
		text := viewNote(update.CallbackQuery.Data)
		_, err := bot.EditMessage(update.CallbackQuery.Message, text)
		if err != nil {
			log.Println(err)
		}
	case "delete":
		text := "Press Ok to Confirm. Cancel to, you know, cancel."
		bot.AddButton("OK", "deleteConfirm-"+strings.Split(update.CallbackQuery.Data, "-")[1])
		bot.AddButton("Cancel", "mainMenu")
		_, err := bot.EditMessage(update.CallbackQuery.Message, text)
		if err != nil {
			log.Println(err)
		}
	case "deleteConfirm":
		text := deleteNote(update.CallbackQuery.Data)
		_, err := bot.EditMessage(update.CallbackQuery.Message, text)
		if err != nil {
			log.Println(err)
		}
	}
}

func mainMenu(update goTel.Update) {
	greetings := []string{"Bonjour", "Hola", "Konnichiwa", "Hey", "Hello"}
	rand.Seed(time.Now().Unix())
	text := greetings[rand.Intn(len(greetings))] + ", "
	bot.AddButton("Create Note", "addNote")
	bot.AddButton("List Notes", "listNotes-1")
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
		_, err := bot.EditMessage(update.CallbackQuery.Message, text)
		if err != nil {
			log.Println("error when sending welcome message", err)
		}
	}
}

func handleNoteQuestions(update goTel.Update, currentNote *pendingNotes) {
	switch currentNote.CurrentStage {
	case 0:
		currentNote.Data.Title = update.Message.Text
		err := bot.DeleteMessage(update.Message)
		if err != nil {
			log.Println(err)
			//bot.AddButton("Menu", "mainMenu")
			//_, _ = bot.EditMessage(currentNote.Message, "And Error Occurred. Try again.")
			//delete(notesList, update.Message.From.ID)
			return
		}
		text := "Alright, Got it. You can now type the content of your note. Please send as a single message."
		currentNote.Message, err = bot.EditMessage(currentNote.Message, text)
		currentNote.CurrentStage++
		if err != nil {
			log.Println(err)
			return
		}
	case 1:
		currentNote.Data.Body = update.Message.Text
		err := bot.DeleteMessage(update.Message)
		if err != nil {
			log.Println(err)
			return
		}
		text := "What's the title of the publication, or Video, or Article?"
		currentNote.Message, err = bot.EditMessage(currentNote.Message, text)
		currentNote.CurrentStage++
		if err != nil {
			log.Println(err)
			return
		}
	case 2:
		currentNote.Data.Publication = update.Message.Text
		err := bot.DeleteMessage(update.Message)
		if err != nil {
			log.Println(err)
			return
		}
		text := "What's the Category of your publications? Please enter a number as indicated.\n"
		for i := 0; i < len(pubCategories); i++ {
			text += fmt.Sprintf("%d: %s\n", i+1, pubCategories[i])
		}
		currentNote.Message, err = bot.EditMessage(currentNote.Message, text)
		currentNote.CurrentStage++
		if err != nil {
			log.Println(err)
			return
		}
	case 3:
		catInt, err := strconv.Atoi(update.Message.Text)
		if err != nil || catInt > len(pubCategories) || catInt <= 0 {
			log.Println(err)
			text := "Please enter a valid number. \nWhat's the Category of your publications? Please enter a number as indicated.\n"
			for i := 0; i < len(pubCategories); i++ {
				text += fmt.Sprintf("%d: %s\n", i+1, pubCategories[i])
			}
			err = bot.DeleteMessage(update.Message)
			if err != nil {
				log.Println(err)
			}
			currentNote.Message, err = bot.EditMessage(currentNote.Message, text)
			if err != nil {
				log.Println(err)
			}
			return
		}
		currentNote.Data.Category = pubCategories[catInt-1]
		err = bot.DeleteMessage(update.Message)
		if err != nil {
			log.Println(err)
			return
		}
		text := "Please enter the tags you want to add to this note. Separate multiple tags with a comma."
		currentNote.Message, err = bot.EditMessage(currentNote.Message, text)
		currentNote.CurrentStage++
		if err != nil {
			log.Println(err)
			return
		}
	case 4:
		//var allTags []tag
		details := update.Message.Text
		err := bot.DeleteMessage(update.Message)
		if err != nil {
			log.Println(err)
		}
		tempTags := strings.Split(details, ",")
		for i := 0; i < len(tempTags); i++ {
			newTag := tag{
				Name:      strings.TrimSpace(tempTags[i]),
				CreatedAt: time.Now(),
			}

			currentNote.Data.Tags = append(currentNote.Data.Tags, newTag)
		}
		//currentNote.Data.Tags = append(currentNote.Data.Tags, allTags
		text := "Alright then. All done. Note to be saved:\n\n"
		text += fmt.Sprintf("Title: %s.\n\nPress 'OK' to continue.", currentNote.Data.Title)
		bot.AddButton("OK", "addNoteOk")
		bot.AddButton("Cancel", "bail")
		bot.MakeKeyboard(1)
		currentNote.Message, err = bot.EditMessage(currentNote.Message, text)
		currentNote.CurrentStage++
		if err != nil {
			log.Println(err)
		}
	default:
		err := bot.DeleteMessage(update.Message)
		if err != nil {
			log.Println(err)
		}
		return
	}
}
