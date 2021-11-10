package main

import (
	"fmt"
	"github.com/yoruba-codigy/goTelegram"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
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


func listAllNotes() string {
	var notes []studyNote
	db.Find(&notes)

	var text string
	for i, note := range notes {
		text += fmt.Sprintf("\n%d. %s.\n", i+1, note.Title)
		log.Println(note.Tags)
		for _, t := range note.Tags {
			text += fmt.Sprintf("#%s ", t.Name)
		}
	}

	return text
}
