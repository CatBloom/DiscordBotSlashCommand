package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	command = []*discordgo.ApplicationCommand{
		{
			Name:        "hello-world",
			Description: "Showcase of a basic slash 'HelloWorld'command",
		},
		{
			Name:        "hi-world",
			Description: "Showcase of a basic slash 'Hi World'command ",
		},
	}
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func main() {
	discord, err := discordgo.New(os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		fmt.Println("Error logging in")
		fmt.Println(err)
	}

	discord.AddHandler(commandHandler)

	err = discord.Open()
	if err != nil {
		fmt.Println(err)
	}

	log.Println("Adding commands...")

	_, err = discord.ApplicationCommandBulkOverwrite(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_GUILD_ID"), command)
	if err != nil {
		fmt.Println(err)
	}

	defer discord.Close()

	fmt.Println("Listening...")
	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, os.Interrupt)
	<-stopBot
}

func commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	switch data.Name {
	case "hello-world":
		err := s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "HelloWorld",
				},
			},
		)
		if err != nil {
			fmt.Println(err)
		}
	case "hi-world":
		err := s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			},
		)
		if err != nil {
			fmt.Println(err)
		}

		//時間のかかる処理を想定
		time.Sleep(time.Second * 5)

		//MessageEmbedを作成
		embed := []*discordgo.MessageEmbed{
			{
				Title:       "HiWorld",
				Description: "Sample HiWorld",
				Type:        "rich",
				Color:       0x7eff00,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Hi",
						Value:  "sample",
						Inline: true,
					},
					{
						Name:   "World",
						Value:  "value",
						Inline: true,
					},
				},
			},
		}

		//Interactionを作成してからEditすると制限にかからない
		//HTTP 404 Not Found, {"message": "Unknown interaction", "code": 10062}を回避できる
		_, err = s.InteractionResponseEdit(
			i.Interaction,
			&discordgo.WebhookEdit{
				Embeds: &embed,
			},
		)
		if err != nil {
			fmt.Println(err)
		}
	}
}
