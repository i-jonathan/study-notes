package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"os"
	"strconv"
	"strings"
)

func initDatabase() *gorm.DB {
	connectionLink := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(connectionLink), &gorm.Config{})
	if err != nil {
		log.Println("Can't connect to db")
		log.Fatalln(err)
		return nil
	}
	err = db.AutoMigrate(&studyNote{}, &tag{})
	if err != nil {
		log.Println("error with auto migration")
		log.Fatalln(err)
		return nil
	}

	return db
}

func createNote(note studyNote) bool {
	err := db.Create(&note).Error
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func listAllNotes(callBackData string, userId int) string {
	var notes []studyNote
	page, _ := strconv.Atoi(strings.Split(callBackData, "-")[1])
	db.Preload(clause.Associations).Scopes(paginate(page)).Find(&notes, "user_id = ?", userId)

	if len(notes) < 1 {
		bot.AddButton("Menu", "mainMenu")
		bot.MakeKeyboard(1)
		return "No notes found."
	}

	var text string
	for i, note := range notes {
		text += fmt.Sprintf("\n\n%d. %s.\n", i+1, note.Title)
		for _, t := range note.Tags {
			text += fmt.Sprintf("#%s  ", t.Name)
		}
		bot.AddButton(strconv.Itoa(i+1), "note-"+strconv.FormatInt(note.ID, 10))
	}
	bot.MakeKeyboard(len(notes))

	var tempNotes []studyNote
	var count int64
	db.Find(&tempNotes, "user_id = ?", userId).Count(&count)

	col := 0
	if page > 1 {
		bot.AddButton("Prev", "listNotes-"+strconv.Itoa(page-1))
		col += 1
	}
	if int64(page * 8) < count {
		bot.AddButton("Next", "listNotes-"+strconv.Itoa(page+1))
		col += 1
	}
	if col != 0 {
		bot.MakeKeyboard(2)
	}

	bot.AddButton("Menu", "mainMenu")
	bot.MakeKeyboard(1)
	return text
}

func paginate(page int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageSize := 8
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func viewNote(callBackData string) string {
	var note studyNote
	noteId, err := strconv.Atoi(strings.Split(callBackData, "-")[1])
	if err != nil {
		log.Println(err)
		return "An error occurred. Please try again later."
	}

	db.Preload(clause.Associations).Find(&note, "id = ?", noteId)
	text := fmt.Sprintf("Title: %s\n\nBody: %s \n\nCategory: %s\n\nPublication: %s\n\n",
		note.Title, note.Body, note.Category, note.Publication)

	for _, t := range note.Tags {
		text += fmt.Sprintf("#%s  ", t.Name)
	}
	bot.AddButton("Back", "listNotes-1")
	bot.MakeKeyboard(1)
	bot.AddButton("Edit", "edit-"+strconv.FormatInt(note.ID, 10))
	bot.AddButton("Delete", "delete-"+strconv.FormatInt(note.ID, 10))
	bot.MakeKeyboard(2)
	bot.AddButton("Menu", "mainMenu")
	return text
}

func deleteNote(callBackData string) string {
	log.Println("test 2")
	noteId, err := strconv.Atoi(strings.Split(callBackData, "-")[1])
	if err != nil {
		log.Println(err)
		return "An error Occurred."
	}
	var tags []tag

	db.Find(&tags, "note_id = ?", noteId)
	for _, t := range tags {
		db.Delete(&t)
	}

	db.Delete(&studyNote{}, noteId)
	bot.AddButton("Menu", "mainMenu")
	bot.MakeKeyboard(1)
	return "Note Deleted"
}

func listTags(userId int) string {
	var tags []tag
	db.Find(&tags, "user_id = ?", userId)

	if len(tags) < 1 {
		return "No tags found"
	}

	tagNames := make(map[string]bool)
	var text string

	for _, t := range tags {
		log.Println(t.Name)
		if _, val := tagNames[t.Name]; !val {
			tagNames[t.Name] = true
			text += "#" + t.Name + "  "
		}
	}

	return text
}

func listNoteByTag(searchData *pendingSearch, userId int, callBackData string) string {
	page, err := strconv.Atoi(strings.Split(callBackData, "-")[1])
	if err != nil {
		log.Println(err)
	}
	tagNames := strings.Split(searchData.Query, ",")
	for i, t := range tagNames {
		tagNames[i] = strings.Title(t)
	}
	var notes []studyNote

	db.Joins("JOIN note_tags ON study_notes.id = note_tags.study_note_id").Joins(
		"JOIN tags on note_tags.tag_id = tags.id and study_notes.user_id=tags.user_id").Where(
			"tags.name in ?", tagNames).Where("tags.user_id = ?", userId).Select(
				"study_notes.id", "title", "publication", "body", "category",
				"study_notes.user_id").Find(&notes)

	if len(notes) < 1 {
		bot.AddButton("Menu", "mainMenu")
		bot.MakeKeyboard(1)
		return "No notes Found."
	}
	var text string

	for i, n := range notes {
		text += fmt.Sprintf("%d. %s\n", i+1, n.Title)

		for _, t := range n.Tags {
			text += fmt.Sprintf("#%s  ", t.Name)
		}
		bot.AddButton(strconv.Itoa(i+1), "note-"+strconv.FormatInt(n.ID, 10))
	}
	bot.MakeKeyboard(len(notes))

	var tempNotes []studyNote
	var count int64
	db.Joins("JOIN note_tags ON study_notes.id = note_tags.study_note_id").Joins(
		"JOIN tags on note_tags.tag_id = tags.id and study_notes.user_id=tags.user_id").Where(
		"tags.name in ?", tagNames).Where("tags.user_id = ?", userId).Select(
		"study_notes.id", "title", "publication", "body", "category",
		"study_notes.user_id").Find(&tempNotes).Count(&count)

	col := 0
	if page > 1 {
		bot.AddButton("Prev", "listNotes-"+strconv.Itoa(page-1))
		col += 1
	}
	if int64(page * 8) < count {
		bot.AddButton("Next", "listNotes-"+strconv.Itoa(page+1))
		col += 1
	}
	if col != 0 {
		bot.MakeKeyboard(2)
	}

	bot.AddButton("Menu", "mainMenu")
	bot.MakeKeyboard(1)
	return text
}