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
	CreatedAt time.Time
	UserId 	  int
}

type studyNote struct {
	ID          int64
	Title       string
	Body        string
	Publication string
	Category    string
	Tags        []tag `gorm:"many2many:note_tags"`
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

type pendingSearch struct {
	Message goTelegram.Message
	Query string
	Page int
}