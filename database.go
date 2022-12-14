package main

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var giveDB = make(chan bool)

var DBChan = make(chan *gorm.DB)

var Done = make(chan bool)

type Guilds struct {
	DefaultChannel string
	GID            string `gorm:"primaryKey"`
}

func InitDB() {

	DB, err := gorm.Open(sqlite.Open("selfbot.db"), &gorm.Config{SkipDefaultTransaction: true})

	if err != nil {
		panic(err)
	}
	DB.AutoMigrate(&Guilds{})

	go returnDB(DB)

}

func returnDB(DB *gorm.DB) {
	for {

		<-giveDB
		DBChan <- DB
		<-Done

	}
}
