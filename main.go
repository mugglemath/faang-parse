package main

import (
	"faang-parse/utils"
	"fmt"
	"log"
)

func main() {
	// load environment variables from .env
	utils.LoadEnvironmentVariables()

	// declare and initialize variables
	port := utils.GetEnv("REDIS_PORT")
	stream := utils.GetEnv("STREAM_NAME")
	group := utils.GetEnv("GROUP_NAME")
	consumer := utils.GetEnv("CONSUMER_NAME")
	apiKey := utils.GetEnv("OPENAI_API_KEY")
	prompt := utils.GetEnv("LLM_PROMPT")

	// get messages from Redis server
	rdb, ctx := utils.CreateRedisClient(port)
	messages, err := utils.ConsumeMessages(stream, group, consumer, ctx, rdb)
	if err != nil {
		log.Fatalf("Error retreiving messages %v", err)
	}

	// send message and get response
	client := utils.CreateOpenAIClient(apiKey)
	for _, msg := range messages {
		chatCompletion, err := utils.SendMessage(client, msg.Content, prompt)
		if err != nil {
			log.Fatalf("Error sending message to API: %v", err)
		}
		fmt.Printf("\n %v \n", chatCompletion.Choices[0].Message.Content)
		fmt.Println("-----------------------------------------------------")
	}
}
