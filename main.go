// This file utilizes a REST API approach.
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
	BLACKHAT = "Consider having a black hat thinking. Black hat thinking is synonymous with being cautious about everything. This thinking is typically useful for preventing anything that could be construed or produced, both intentionally and unintentionally as illegal, immoral, unprofitable, polluting, toxic and similar things that might go against the organisation's values or general public perception. This thinking induces a sense of mistrust in every step the person takes to prevent future troubles. In such a case, the person always assumes the worst possible outcome will happen. It is very much in line with Murphy’s Law - 'Anything that can go wrong will go wrong.'\nSuch a mindset is beneficial in fighting confirmation bias, which lets everyone and the proposer be critical of the idea and find potential pitfalls. This is an exercise in the analytical side of the brain causing the person to be cautious and fault-finding when needed and scrutinise every nook and corner.\nThis type of thinking is useful when there is a need to find fallacies in logic and scan for dangers. Any answer provided by this reason must accompany a reason to justify the logic behind it. The thinking assumes that the proposed plan is a failure at present and is tasked to find out all the causes."

	YELLOWHATHAT = "Consider Yellow Hat thinking. This method of thinking signifies sunshine and optimism in a mindset. This type of thinking usually looks out for opportunities and benefits. In this type of thinking, the user conducts assessment to see more benefits. The thinking supports self-optimism and self-pride with also benefitting from ideas that comes from every stream. The user is always on the lookout for new opportunities and ideas for the idea. There is always a conscious approach to look at the bright side in every situation. There is also a constructive thinking model that helps the user to find solutions in every situation and predicament. Yellow Hat thinking is more in line with Yhprum’s Law which says that, “Anything that can go right, will go right”.\nConsider two scenarios, in first scenario, there is the thought of playing lottery and hoping to win while in the second scenario, there is the idea of putting a man on the moon. The second scenario is more in line with Yellow Hat thinking as it focuses on turning improbable to possibly true. This line of thinking is heavily backed with reasons and logical support as it answers the rudimentary question of “What is the view based upon?”. The question is answered on the basis of experience, available information, logical deductions, trends, guesses and hopes. The thinking works together with reasons and optimism. Speculation plays a big part in this line of thinking. Even with substantial evidence, yellow hat thinking can now put point forward in a speculative manner. The approach to yellow hat thinking is directly associated with constructive thinking as in an assessment scenario, yellow hat thinking strives to find benefits in an idea. Bold thoughts could be put forward with conviction using this line of thought."

	REDHAT = ""
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

	//model settings
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

	// System Instructions
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(YELLOWHATHAT)},
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
