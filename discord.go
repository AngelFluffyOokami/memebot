package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

var s *discordgo.Session
var err error
var Channels []string

var MediaLink = make(chan string)

func initDiscord(token string) *discordgo.Session {

	s, err = discordgo.New(token)
	s.Identify.Intents = discordgo.IntentsAll
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err = s.Open()
	if err != nil {
		panic(err)
	}
	s.AddHandler(OnGuildAdd)
	s.AddHandler(OnGuildRemove)
	s.AddHandler(OnMessageCreate)

	PopulateServers(s)
	PopulateChannels()

	go Update(s)

	return s
}

func PopulateChannels() {
	giveDB <- true
	DB := <-DBChan
	var guilds []Guilds
	DB.Find(&guilds)
	Done <- true
	for x := range guilds {
		if guilds[x].DefaultChannel != "" {
			Channels = append(Channels, guilds[x].DefaultChannel)
		}

	}
}
func Update(s *discordgo.Session) {

	for {
		Message := <-MediaLink

		for x := range Channels {

			s.ChannelMessageSend(Channels[x], Message)

			seconds := 1
			time.Sleep(time.Duration(seconds) * time.Second)
		}

	}
}

func PopulateServers(s *discordgo.Session) {
	giveDB <- true
	DB := <-DBChan
	guilds := s.State.Guilds

	for x := range guilds {
		guild := Guilds{
			GID: guilds[x].ID,
		}
		result := DB.First(&guild)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			DB.Save(&guild)
		}
	}
	Done <- true

}

func OnGuildRemove(s *discordgo.Session, g *discordgo.GuildDelete) {
	giveDB <- true
	DB := <-DBChan

	guild := Guilds{
		GID: g.Guild.ID,
	}

	result := DB.First(&guild)

	for x := range Channels {
		if Channels[x] == guild.DefaultChannel {
			Channels = append(Channels[:x], Channels[x+1:]...)
			break
		}

	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		DB.Delete(&guild)
	}

	fmt.Println("Guild removed")
	Done <- true
}

func OnGuildAdd(s *discordgo.Session, g *discordgo.GuildCreate) {
	giveDB <- true
	DB := <-DBChan
	guild := Guilds{
		GID: g.Guild.ID,
	}
	result := DB.First(&guild)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		DB.Create(&guild)
	}
	fmt.Println("guild added")
	Done <- true
}
