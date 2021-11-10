package main

import (
	"github.com/yoruba-codigy/goTelegram"
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
