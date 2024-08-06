package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

const (
	BLACKHAT = ""
	WHITEHAT = ""
	GREENHAT = ""
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Failed to load .env file")
	}

	port := os.Getenv("PORT")

	router := gin.Default()
	v1 := router.Group("/api/v1")

	{
		v1.POST("/generate", getResults)
	}

	router.Run(port)
}

func getResults(c *gin.Context) {
	ctx := context.Background()
	var responseTest struct {
		Prompt string `json:"prompt"`
	}

	if err := c.BindJSON(&responseTest); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	apiKey := os.Getenv("API_KEY")

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	model.SetTemperature(1.90)
	model.SetTopK(64)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(4096)
	model.ResponseMIMEType = "text/plain"
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockOnlyHigh,
		},
	}

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(BLACKHAT)},
	}

	session := model.StartChat()
	session.History = []*genai.Content{}

	resp, err := session.SendMessage(ctx, genai.Text(responseTest.Prompt))

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var response genai.Part

	if resp != nil {
		candidates := resp.Candidates
		if candidates != nil {
			for _, candidate := range candidates {
				content := candidate.Content
				if content != nil {
					response = content.Parts[0]
				}
			}
		} else {
			log.Printf("Candidates is nil.\n")
			log.Print(resp.PromptFeedback.BlockReason.String())
		}
	}

	c.JSON(200, response)
}

// func prompts() {
// 	apikey := os.Getenv("API_KEY")
// 	ctx := context.Background()

// 	client, err := genai.NewClient(ctx, option.WithAPIKey(apikey))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer client.Close()

// 	model := client.GenerativeModel("gemini-1.5-flash")

// 	model.SetTemperature(2)
// 	model.SetTopK(64)
// 	model.SetTopP(0.95)
// 	model.SetMaxOutputTokens(4096)
// 	model.ResponseMIMEType = "text/plain"
// 	model.SafetySettings = []*genai.SafetySetting{
// 		{
// 			Category:  genai.HarmCategoryHarassment,
// 			Threshold: genai.HarmBlockOnlyHigh,
// 		},
// 		{
// 			Category:  genai.HarmCategoryHateSpeech,
// 			Threshold: genai.HarmBlockOnlyHigh,
// 		},
// 		{
// 			Category:  genai.HarmCategorySexuallyExplicit,
// 			Threshold: genai.HarmBlockOnlyHigh,
// 		},
// 		{
// 			Category:  genai.HarmCategoryDangerousContent,
// 			Threshold: genai.HarmBlockOnlyHigh,
// 		},
// 	}

// 	model.SystemInstruction = &genai.Content{
// 		Parts: []genai.Part{
// 			genai.Text("You are a melodic rapper with a laid back vibe. You draw inspiration from J.Cole, Saba, Lupe Fiasco, Smino, Baby Keem, A Tribe Called Quest and Run DMC to name a few."),
// 		},
// 	}

// 	resp, err := model.GenerateContent(ctx, genai.Text("write a small jingle of 8 lines about ice cream trucks on a summer day"))

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if resp != nil {
// 		candidates := resp.Candidates
// 		if candidates != nil {
// 			for _, candidate := range candidates {
// 				content := candidate.Content
// 				if content != nil {
// 					text := content.Parts[0]
// 					role := content.Role
// 					log.Printf("%s: %s \n", role, text)
// 				}
// 			}
// 		} else {
// 			log.Printf("Candidates is nil.\n")
// 			log.Print(resp.PromptFeedback.BlockReason.String())
// 		}
// 	}
// }
