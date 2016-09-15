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
			slackBot.sendMessage(msg)
			continue
		}

		// NOTE: the Message object is copied, this is intentional
		go func(msg Message) {

			intent, err := response.Intent()
			if err != nil {
				// Handle error
				msg.Text = fmt.Sprintf("Intent error: %s", err)
				slackBot.sendMessage(msg)
				return
			}

			msg.Text = fmt.Sprintf("Intent found: %s", intent)
			slackBot.sendMessage(msg)

			allEntities := response.AllEntities()
			for key, value := range allEntities {
				for _, entity := range value {
					msg.Text = fmt.Sprintf("Entity[%v].name=%s\nEntity[%v].raw=%s\nEntity[%v].formated=%v\n",
						key, entity.Name(), key, entity.Raw(), key, entity.Field("formated"))
					slackBot.sendMessage(msg)
				}
			}

			msg.Text = fmt.Sprintf("response language: %s\nresponse status: %v\nresponse timestamp: %s\nresponse version: %s\n",
				response.Language(), response.Status(), response.Timestamp(), response.Version())
			slackBot.sendMessage(msg)

			sentence := response.Sentence()

			msg.Text = fmt.Sprintf("sentence.Source()=%s\nsentence.Type()=%s\nsentence.Action()=%s\nsentence.Agent()=%s\nsentence.Polarity()=%s\n",
				sentence.Source(), sentence.Type(), sentence.Action(), sentence.Agent(), sentence.Polarity())
			slackBot.sendMessage(msg)
		}(msg)

	}
}
