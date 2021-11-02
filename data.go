package main

import (
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
type tags struct {
	Id        int64
	Name      string
	CreatedAt time.Time
}

type studyNotes struct {
	Id          int64
	Title       string
	Body        string
	Publication string
	Category    string
	Tags        []tags
	UserId      int
	CreatedAt   time.Time
}

// structure for processing notes
type pendingNotes struct {
	Stages       int
	CurrentStage int
	Message      goTelegram.Message
	Data         studyNotes
}

func initDatabase() *gorm.DB {
	connectionLink := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(connectionLink), &gorm.Config{})
	if err != nil {
		log.Println("Can't connect to db")
		log.Fatalln(err)
		return nil
	}
	err = db.AutoMigrate(&tags{}, &studyNotes{})
	if err != nil {
		log.Println("error with auto migration")
		log.Fatalln(err)
		return nil
	}

	return db
}

func createNote(note studyNotes) {
	db.Create(note)
}

//
//func listAllNotes() text {
//
//}
