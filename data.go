package main

import (
	"fmt"
	"github.com/yoruba-codigy/goTelegram"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

//type category struct {
//	ID int64
//	Name string
//	CreatedAt time.Time
//}

// database tables
type tag struct {
	ID        int64
	Name      string
	NoteID    int
	CreatedAt time.Time
}

type studyNote struct {
	ID          int64
	Title       string
	Body        string
	Publication string
	Category    string
	Tags        []tag `gorm:"foreignKey:NoteID"`
	UserId      int
	CreatedAt   time.Time
}

// structure for processing notes
type pendingNotes struct {
	Stages       int
	CurrentStage int
	Message      goTelegram.Message
	Data         studyNote
}

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


func listAllNotes(callBackData string) string {
	var notes []studyNote
	page, _ := strconv.Atoi(strings.Split(callBackData, "-")[1])
	db.Preload(clause.Associations).Scopes(paginate(page)).Find(&notes)

	var text string
	for i, note := range notes {
		text += fmt.Sprintf("\n\n%d. %s.\n\n", i+1, note.Title)
		for _, t := range note.Tags {
			text += fmt.Sprintf("#%s  ", t.Name)
		}
		bot.AddButton(strconv.Itoa(i+1), "note-"+strconv.FormatInt(note.ID, 10))
	}
	bot.MakeKeyboard(len(notes))

	var tempNotes []studyNote
	var count int64
	db.Find(&tempNotes).Count(&count)

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