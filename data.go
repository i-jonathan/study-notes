package main

import (
	"github.com/yoruba-codigy/goTelegram"
	"time"
)

type category struct {
	Id int64
	Name string
	CreatedAt time.Time
}

type tags struct {
	Id int64
	Name string
	CreatedAt time.Time
}

type studyNotes struct {
	Id int64
	Title string
	Body string
	Publication string
	Category category
	Tags []tags
	UserId int
	CreatedAt time.Time
}

type pendingNotes struct {
	Stages int
	CurrentStage int
	Message goTelegram.Message
	Data studyNotes
}

//func createSchema(db *pg.DB) error {
//	models := []interface{}{
//		(*category)(nil),
//		(*tags)(nil),
//		(*studyNotes)(nil),
//	}
//
//	for _, model := range models {
//		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
//			Temp: true,
//		})
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}