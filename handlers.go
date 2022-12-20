package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func DeleteServer(i *discordgo.MessageCreate) {
	giveDB <- true
	DB := <-DBChan

	guild := Guilds{
		GID: i.Message.GuildID,
	}

	result := DB.First(&guild)

	for x := range Channels {
		if Channels[x] == guild.DefaultChannel {
			Channels = append(Channels[:x], Channels[x+1:]...)
			break
		}

	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		guild.DefaultChannel = ""
		DB.Save(&guild)
	}

	Done <- true
	fmt.Println("Server removed")
}

func OnMessageCreate(s *discordgo.Session, i *discordgo.MessageCreate) {
	channel, _ := s.State.Channel(i.Message.ChannelID)

	if channel.Type == 1 {
		fmt.Println("DM received: " + i.Message.Content)

		if i.Message.Author.ID == OwnerID {
			fmt.Println("Message is owner")

			MediaLink <- i.Message.Content
			MediaList := i.Message.Attachments
			for x := range MediaList {
				MediaLink <- MediaList[x].ProxyURL
				fmt.Println(MediaList[x].ProxyURL)
			}

		}
	} else if i.Message.Author.ID == OwnerID {
		if strings.ToLower(i.Message.Content) == "repost memes here" {
			UpdateChannel(i)
			fmt.Println("Channel Set: " + i.Message.ChannelID)
		} else if strings.ToLower(i.Message.Content) == "stop sending memes here" {
			DeleteServer(i)
		}
	}

}

func UpdateChannel(i *discordgo.MessageCreate) {
	giveDB <- true
	DB := <-DBChan

	guild := Guilds{
		GID: i.Message.GuildID,
	}

	oldGuild := Guilds{
		GID: i.Message.GuildID,
	}
	DB.First(&guild)
	DB.First(&oldGuild)
	guild.DefaultChannel = i.Message.ChannelID

	DB.Save(&guild)
	Done <- true
	Channels = append(Channels, i.Message.ChannelID)
	if Channels != nil {
		for x := range Channels {
			if Channels[x] == oldGuild.DefaultChannel {
				if len(Channels) >= 2 {
					Channels = append(Channels[:x], Channels[x+1:]...)
					break
				} else {
					Channels = nil
				}

			}
		}
	} else {
		Channels = append(Channels, i.Message.GuildID)
	}

}
