package utils

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

// creates an OpenAI client
func CreateOpenAIClient(apiKey string) *openai.Client {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	return client
}

// sends a message to OpenAI API and returns the response
func SendMessage(client *openai.Client, content string, prompt string) (data *openai.ChatCompletion, err error) {
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{openai.ChatCompletionUserMessageParam{
			Role:    openai.F(openai.ChatCompletionUserMessageParamRoleUser),
			Content: openai.F[openai.ChatCompletionUserMessageParamContentUnion](shared.UnionString(prompt + " " + content)),
		}}),
		Model: openai.F(openai.ChatModelGPT4oMini),
	})
	if err != nil {
		panic(err.Error())
	}
	return chatCompletion, nil
}
