package main

import (
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

	botToken := viper.GetString("token.RecastAI")

	client := recast.NewClient(botToken, "en")

	text := "I need to know the weather in Paris"
	response, err := client.TextRequest(text, nil)
	if err != nil {
		// Handle error
		log.Println("TextRequest error:", err)
		return
	}

	intent, err := response.Intent()
	log.Println("Intent found:", intent)

	allEntities := response.AllEntities()
	for key, value := range allEntities {
		log.Printf("Entity[%v]\n", key)
		for _, entity := range value {
			log.Printf("          .name=%s\n", entity.Name())
			log.Printf("          .raw=%s\n", entity.Raw())
			log.Printf("          .formated=%v\n", entity.Field("formated"))
		}
	}

	log.Println("response languages:", response.Language())
	log.Println("response status:", response.Status())
	log.Println("response timestamp:", response.Timestamp())
	log.Println("response version:", response.Version())

	sentence := response.Sentence()
	log.Println("sentence.Source()=", sentence.Source())
	log.Println("sentence.Type()=", sentence.Type())
	log.Println("sentence.Action()=", sentence.Action())
	log.Println("sentence.Agent()=", sentence.Agent())
	log.Println("sentence.Polarity()=", sentence.Polarity())
}
