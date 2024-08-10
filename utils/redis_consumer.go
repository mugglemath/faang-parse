package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type Message struct {
	MessageId  string
	JobId      string
	Title      string
	DatePosted string
	Company    string
	Content    string
}

// create new Redis client
func CreateRedisClient(port string) (*redis.Client, context.Context) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:" + port,
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	// debug
	if DebugEnabled() {
		fmt.Printf("Connected to Redis on %v\n", rdb.Options().Addr)
	}
	return rdb, ctx
}

// consume Redis stream messages from server
func ConsumeMessages(
	streamName string,
	groupName string,
	consumerName string,
	ctx context.Context,
	client *redis.Client,
) ([]Message, error) {

	// fast fail if stream not found or nil
	streamLength, err := client.XLen(ctx, streamName).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting stream length: %v", err)
	}

	if streamLength == 0 {
		fmt.Println("There are no entries in the stream.")
		return []Message{}, nil
	}

	// debug
	if DebugEnabled() {
		fmt.Printf("Number of entries in the stream: %v\n", streamLength)
	}

	// fast fail if messages already consumed
	groupInfo, err := client.XInfoGroups(ctx, streamName).Result()
	if err != nil {
		log.Fatalf("Error getting group info: %v", err)
	}

	for _, group := range groupInfo {
		if group.EntriesRead == streamLength {
			fmt.Println("No new messages to read")
			return []Message{}, nil
		}
	}

	// read messages from stream
	var messages []Message
	streamMessages, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    groupName,
		Consumer: consumerName,
		Streams:  []string{streamName, ">"},
		Count:    streamLength,
		Block:    1,
	}).Result()

	if err != nil {
		if err == redis.Nil {
			fmt.Println("No new messages to read")
			return messages, nil
		}
		fmt.Printf("Error reading from stream: %v\n", err)
		return messages, err
	}

	// process messages
	for _, msg := range streamMessages {
		for _, message := range msg.Messages {
			newMessage := Message{
				MessageId:  message.ID,
				JobId:      message.Values["jobId"].(string),
				Title:      message.Values["title"].(string),
				DatePosted: message.Values["datePosted"].(string),
				Company:    message.Values["company"].(string),
				Content:    message.Values["content"].(string),
			}
			// append message to slice
			messages = append(messages, newMessage)

			// ack message
			_, err := client.XAck(ctx, streamName, groupName, message.ID).Result()
			if err != nil {
				return nil, fmt.Errorf("error acknowledging message: %v", err)
			}
		}
	}
	// debug
	if DebugEnabled() {
		for _, msg := range messages {
			fmt.Printf("Message ID: %s\n", msg.MessageId)
			fmt.Printf("Job ID: %s\n", msg.JobId)
			fmt.Printf("Title: %s\n", msg.Title)
			fmt.Printf("Date Posted: %s\n", msg.DatePosted)
			fmt.Printf("Company: %s\n", msg.Company)
			fmt.Printf("Content: %s\n", msg.Content)
			fmt.Println("-----------------------------------------------------")
		}
	}
	return messages, nil
}
