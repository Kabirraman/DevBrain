package chat

import (
	"context"

	"github.com/Kabirraman/DevBrain/internal/ai"

	"google.golang.org/genai"
)

func GenerateAnswer(
	graphContext string,
	question string,
) (string, error) {

	client, err := ai.NewGeminiClient()

	if err != nil {
		return "", err
	}

	prompt := `
You are DevBrain.

Answer ONLY using the provided knowledge graph.

If the answer is not present in the graph,
say:

"I couldn't find enough information in the current knowledge graph."

Knowledge Graph:

` + graphContext + `

Question:

` + question

	resp, err := client.Models.GenerateContent(
		context.Background(),
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)

	if err != nil {
		return "", err
	}

	return resp.Text(), nil
}