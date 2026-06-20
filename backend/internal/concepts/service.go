package concepts

import (
	"strings"
	"fmt"
	"context"
	"encoding/json"
	"os"
	"github.com/Kabirraman/DevBrain/internal/database"
	"github.com/Kabirraman/DevBrain/internal/models"
	"google.golang.org/genai"
)


type ConceptResponse struct {
	Concepts []string `json:"concepts"`
}

func ExtractConcepts(content string) ([]string, error) {

	client, err := genai.NewClient(
		context.Background(),
		&genai.ClientConfig{
			APIKey: os.Getenv("GEMINI_API_KEY"),
		},
	)

	if err != nil {
		return nil, err
	}

	prompt := `
Extract technical concepts from the text.

Return ONLY valid JSON.

Do NOT use markdown.
Do NOT wrap the response in triple backticks.
Do NOT explain anything.

Output format:

{
  "concepts":[]
}

Example:

{
  "concepts":[
    "Go",
    "Goroutines",
    "Concurrency"
  ]
}

Text:
` + content

	resp, err := client.Models.GenerateContent(
		context.Background(),
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)

	if err != nil {
		return nil, err
	}
	fmt.Println("GEMINI RESPONSE:")
	fmt.Println(resp.Text())

	var result ConceptResponse

	cleaned := cleanJSONResponse(
	resp.Text(),
)

err = json.Unmarshal(
	[]byte(cleaned),
	&result,
)

	if err != nil {
		return nil, err
	}

	return result.Concepts, nil
}

func SaveConcepts(
	concepts []string,
) error {

	for _, concept := range concepts {

		c := models.Concept{
			Name: concept,
		}

		err := database.DB.Create(&c).Error

		if err != nil {
			return err
		}
	}

	return nil
}

func cleanJSONResponse(text string) string {

	text = strings.TrimSpace(text)

	text = strings.TrimPrefix(
		text,
		"```json",
	)

	text = strings.TrimPrefix(
		text,
		"```",
	)

	text = strings.TrimSuffix(
		text,
		"```",
	)

	return strings.TrimSpace(text)
}