package main

import (
	"fmt"
	"log"

	"github.com/RecastAI/SDK-Golang/recast"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		log.Println("No configuration file loaded - using defaults")
		return
	}

	recastBotToken := viper.GetString("token.RecastAI")
	slackBotToken := viper.GetString("token.Slack")

	// start a websocket-based Real Time API session
	slackBot, err := NewBot(slackBotToken)
	if err != nil {
		log.Fatalf("Can't connect to Slack [%s]", err)
	}
	fmt.Println("SlackGoBot is running...")
	fmt.Println("SlackGoBot id is =", slackBot.id)

	recastClient := recast.NewClient(recastBotToken, "en")

	for {

		// Read each incoming message
		msg, err := slackBot.receiveMessage()
		if err != nil {
			log.Fatal("Error while getting message", err)
		}

		if msg.Type != "message" {
			continue
		}

		// Test if the message was written by the bot
		if msg.User == slackBot.id {
			continue
		}

		// The received message is a 'message' type.

		response, err := recastClient.TextRequest(msg.Text, nil)
		if err != nil {
			// Handle error
			msg.Text = fmt.Sprintf("TextRequest error: %s", err)
			fmt.Println("msg1=", msg)
			slackBot.sendMessage(msg)
			continue
		}

		// NOTE: the Message object is copied, this is intentional
		go func(msg Message) {

			for intent := range response.Intents {
				msg.Text = fmt.Sprintf("Intent found: %v", intent)
				fmt.Println("msg2=", msg)
				slackBot.sendMessage(msg)
			}

			for key, value := range response.Entities {
				for _, entity := range value {
					msg.Text = fmt.Sprintf("Entity[%v].name=%s\nEntity[%v].confidence=%f\n",
						key, entity.Name, key, entity.Confidence)
					fmt.Println("msg3=", msg)
					slackBot.sendMessage(msg)
				}
			}

			msg.Text = fmt.Sprintf("response UUID: %s\nresponse source: %s\nresponse act: %s\nresponse type: %s\nresponse sentiment: %s\nresponse language: %s\nresponse status: %d\nresponse timestamp: v\nresponse version: %s\n",
				response.UUID, response.Source, response.Act, response.Type, response.Sentiment, response.Language, response.Status, response.Version)
			fmt.Println("msg4=", msg)
			slackBot.sendMessage(msg)

		}(msg)

	}
}
