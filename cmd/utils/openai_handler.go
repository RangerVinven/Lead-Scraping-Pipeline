package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func GenerateAbstracts(markdownPages []string) []string {
	var prompt = `
You're provided a Markdown scrape of a website page. Your task is to provide a two-paragraph abstract of what this page is about.

Return in this JSON format:

{"abstract":"your abstract goes here"}

Rules:
- Your abstract should be comprehensive—similar level of detail as an academic abstract.
- Use a straightforward, spartan tone of voice.
- If the page has no content, return: {"abstract": "no content"}
`

	var abstracts = []string{}

	for _, markdown := range markdownPages {
		client := openai.NewClient(
			option.WithAPIKey(os.Getenv("OPENAI_API_KEY")), // defaults to os.LookupEnv("OPENAI_API_KEY")
		)
		chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(prompt),
				openai.UserMessage(markdown),
			},
			Model: openai.ChatModelGPT4_1Nano,
		})
		if err != nil {
			panic(err.Error())
		}

		var response = chatCompletion.Choices[0].Message.Content

		abstracts = append(abstracts, response)
	}

	return abstracts
}

func GenerateIcebreaker(abstracts []string, email string, companyName string) (string, string) {
	var combinedAbstracts = strings.Join(abstracts, "\n---\n")

	var prompt = fmt.Sprintf(`
We just scraped a series of web pages for a business. Your task is two-fold:
1. Identify the First Name of the most appropriate person to contact using the provided email: %s and Company Name: %s
2. Use that information to write a personalized cold email icebreaker.

Name Extraction Rules (In order of priority):
- FIRST: Cross-reference the email prefix (the part before @) with the website summaries to find a full name match.
- SECOND: If the email is generic (info, office, etc.) but the Company Name contains a recognizable human first name (e.g., "Stuart's Plumbing", "David Wood Heating"), use that name.
- THIRD: If the email prefix itself is a full name like 'sara@', use 'Sara'.
- NEVER return a single initial (e.g., 'D').
- NEVER return system names (e.g., 'Info', 'Admin', 'Office').
- If no name can be found with high confidence, return "Unknown".

Icebreaker Rules:
- If name is "Unknown", start with "Hey,". Otherwise, start with "Hey {FirstName},".
- Use a spartan/laconic tone.
- Shorten company and location names when possible.
- Avoid obvious compliments. Focus on small, unique details from the summaries.
- Talk in first and second person only ("I" and "you").

Return your response in this JSON format:
{"firstName": "Name or Unknown", "icebreaker": "Your personalized message here"}
`, email, companyName)
	
	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
	)
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompt),
			openai.UserMessage(combinedAbstracts),
		},
		Model: openai.ChatModelGPT4_1Nano,
	})
	if err != nil {
		panic(err.Error())
	}

	var jsonStr = chatCompletion.Choices[0].Message.Content
	
	var data map[string]string
	err = json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		fmt.Println("Couldn't parse the JSON. Received error:", err)
		return "", ""
	}

	name := data["firstName"]
	if name == "Unknown" || name == "unknown" {
		name = ""
	}

	return name, data["icebreaker"]
}
