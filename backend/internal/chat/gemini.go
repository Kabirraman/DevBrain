package chat

import (
	"context"
	"fmt"

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
You are DevBrain, a friendly technical mentor.

Answer the learner's question by explaining it in your own words, as a
short, clear paragraph - the way a mentor would talk to a student, not
by listing raw facts.

Use the knowledge graph below as your only source of truth. Weave the
relevant facts into normal sentences. Do NOT output the graph lines
verbatim (no "X RELATION Y" formatting) and do NOT just list concepts -
write real sentences.

If the graph doesn't contain enough to answer, say exactly:
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
		fmt.Println("GEMINI API ERROR (chat):", err)
		return "", err
	}

	return resp.Text(), nil
}