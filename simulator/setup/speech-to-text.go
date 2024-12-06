package functions

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ermes-labs/api-go/api"
	rc "github.com/ermes-labs/storage-redis/packages/go"
	"github.com/sashabaranov/go-openai"
)

var client = openai.NewClient(os.Getenv("OPENAI_API_KEY"))

func speech_to_text(
	w http.ResponseWriter,
	r *http.Request,
	sessionToken *api.SessionToken,
) (err error) {
	// Read the file from the request body.
	fileWav := r.Body
	defer r.Body.Close()
	if err != nil {
		return err
	}

	// Create a session if it does not exists and acquire it.
	return Node.CreateAndAcquireSession(
		r.Context(),
		&sessionToken,
		api.CreateAndAcquireSessionOptions{},
		func(sessionToken api.SessionToken) error {
			ks := rc.NewErmesKeySpaces(sessionToken.SessionId)
			// Code to handle speech to text and AIChat.
			prompt := SpeechToText(fileWav)
			aiResponse := Ask(prompt)

			// Derive the CPU usage from the response size. Early returns and
			// unset resources will default to the average(!) usage of 1.
			Node.SetResourcesTotalUsage(r.Context(), sessionToken, map[string]float64{
				"cpu": deriveCPUUsage(len(aiResponse)),
			})

			// If there is an error it will be lifted up to the main scope.
			err = redisClient.RPush(r.Context(), ks.Session("chat"), aiResponse).Err()
			if err != nil {
				return err
			}

			// Write the response.
			w.Write([]byte(aiResponse))
			return nil
		})
}

func SpeechToText(fileReader io.Reader) string {
	// TODO: This part is commented out to remove the little uncertainty given by
	// The use of an external API. Initial tests are done in a more predictable way
	// To ease the comparison between the different setups.
	/*
		ctx := context.Background()
		audioRequest := openai.AudioRequest{
			Model:  "whisper-1",
			Reader: fileReader,
		}

		response, err := client.CreateTranscription(ctx, audioRequest)
		if err != nil {
			log.Printf("Failed to transcribe audio: %v", err)
			return ""
		}

		return (response.Text)
	*/

	stringFile, err := io.ReadAll(fileReader)
	if err != nil {
		panic(err)
	}

	start := time.Now()
	for time.Since(start) < time.Duration((50+int(len(stringFile)/1024/100)))*time.Millisecond {
		// Keep the CPU busy
	}

	return string(stringFile)
}

func Ask(prompt string) string {
	// TODO: This part is commented out to remove the little uncertainty given by
	// The use of an external API. Initial tests are done in a more predictable way
	// To ease the comparison between the different setups.
	/*ctx := context.Background()

	req := openai.CompletionRequest{
		Model:     openai.GPT3Ada,
		MaxTokens: 50,
		Prompt:    prompt,
	}
	resp, err := client.CreateCompletion(ctx, req)
	if err != nil {
		fmt.Printf("Completion error: %v\n", err)
		return ""
	}
	return resp.Choices[0].Text
	*/

	start := time.Now()
	for time.Since(start) < time.Duration((50+int(len(prompt)/1024/100)))*time.Millisecond {
		// Keep the CPU busy
	}
	response := prompt
	return response
}

func deriveCPUUsage(size int) float64 {
	// Derive the CPU usage from the response size. Early returns and unset
	// resources will default to the average(!) usage of 1.
	return 0.03 * (1 + float64(size)/1024/5)
}
