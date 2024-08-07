package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

const (
	BLACKHAT = "Consider having a black hat thinking. Black hat thinking is synonymous with being cautious about everything. This thinking is typically useful for preventing anything that could be construed or produced, both intentionally and unintentionally as illegal, immoral, unprofitable, polluting, toxic and similar things that might go against the organisation's values or general public perception. This thinking induces a sense of mistrust in every step the person takes to prevent future troubles. In such a case, the person always assumes the worst possible outcome will happen. It is very much in line with Murphy’s Law - 'Anything that can go wrong will go wrong.'\nSuch a mindset is beneficial in fighting confirmation bias, which lets everyone and the proposer be critical of the idea and find potential pitfalls. This is an exercise in the analytical side of the brain causing the person to be cautious and fault-finding when needed and scrutinise every nook and corner.\nThis type of thinking is useful when there is a need to find fallacies in logic and scan for dangers. Any answer provided by this reason must accompany a reason to justify the logic behind it. The thinking assumes that the proposed plan is a failure at present and is tasked to find out all the causes."

	YELLOWHAT = "Consider Yellow Hat thinking. This method of thinking signifies sunshine and optimism in a mindset. This type of thinking usually looks out for opportunities and benefits. In this type of thinking, the user conducts assessment to see more benefits. The thinking supports self-optimism and self-pride with also benefitting from ideas that comes from every stream. The user is always on the lookout for new opportunities and ideas for the idea. There is always a conscious approach to look at the bright side in every situation. There is also a constructive thinking model that helps the user to find solutions in every situation and predicament. Yellow Hat thinking is more in line with Yhprum’s Law which says that, “Anything that can go right, will go right”.\nConsider two scenarios, in first scenario, there is the thought of playing lottery and hoping to win while in the second scenario, there is the idea of putting a man on the moon. The second scenario is more in line with Yellow Hat thinking as it focuses on turning improbable to possibly true. This line of thinking is heavily backed with reasons and logical support as it answers the rudimentary question of “What is the view based upon?”. The question is answered on the basis of experience, available information, logical deductions, trends, guesses and hopes. The thinking works together with reasons and optimism. Speculation plays a big part in this line of thinking. Even with substantial evidence, yellow hat thinking can now put point forward in a speculative manner. The approach to yellow hat thinking is directly associated with constructive thinking as in an assessment scenario, yellow hat thinking strives to find benefits in an idea. Bold thoughts could be put forward with conviction using this line of thought."

	REDHAT = ""
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading environment file")
	}

	apiKey := os.Getenv("API_KEY")

	if apiKey == "" {
		log.Fatal("No API_KEY environment variable set.")
		return
	}

	titleFigure := figure.NewColorFigure("HAT PROMPTS", "", "green", true)
	titleFigure.Print()

	fmt.Print("Type exit to exit \n\n")

	chatWithLLM(apiKey)
}

func chatWithLLM(apiKey string) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))

	if err != nil {
		log.Fatal(err.Error())
	}

	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	model.SetTemperature(1.5)
	model.SetTopK(64)
	model.SetTopP(1)
	model.SetMaxOutputTokens(2048)
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
		Parts: []genai.Part{genai.Text(YELLOWHAT)},
	}

	session := model.StartChat()
	session.History = []*genai.Content{}

	preresp, err := session.SendMessage(ctx, genai.Text("Provide an introduction about your type of thinking"))

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	if preresp != nil {
		candidates := preresp.Candidates
		if candidates != nil {
			for _, candidate := range candidates {
				content := candidate.Content
				if content != nil {
					log.Print(content.Parts[0])
				}
			}
		} else {
			log.Println("Empty candidate")
			log.Println(preresp.PromptFeedback.BlockReason.String())
		}
	}

	var usermsg string
	reader := bufio.NewReader(os.Stdin)

	for {
		log.Printf("Enter Message:> ")
		input, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal(err.Error())
		}

		usermsg = strings.TrimSpace(input)

		if usermsg == "exit" {
			break
		}

		resp, err := session.SendMessage(ctx, genai.Text(usermsg))

		if err != nil {
			log.Fatal(err.Error())
		}

		if resp != nil {
			candidates := resp.Candidates
			if candidates != nil {
				for _, candidate := range candidates {
					content := candidate.Content
					if content != nil {
						log.Print(content.Parts[0])
					}
				}
			} else {
				log.Println("Empty candidate")
				log.Println(resp.PromptFeedback.BlockReason.String())
			}
		}
	}

	fmt.Print("ENDED CHAT")
}
