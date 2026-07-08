package search

import (
	"context"
	"sort"
	"strings"

	"github.com/Kabirraman/DevBrain/internal/ai"
	"github.com/Kabirraman/DevBrain/internal/database"
	"github.com/Kabirraman/DevBrain/internal/models"

	"google.golang.org/genai"
)

var stopWords = map[string]bool{
	"what": true, "is": true, "are": true, "how": true, "does": true,
	"do": true, "the": true, "a": true, "an": true, "to": true,
	"of": true, "in": true, "for": true, "and": true,
}

func Search(userID string, query string) (*SearchResponse, error) {

	var resources []models.Resource

	err := database.DB.
		Where("user_id = ?", userID).
		Find(&resources).
		Error

	if err != nil {
		return nil, err
	}

	words := tokenize(query)

	type scored struct {
		resource models.Resource
		score    int
	}

	var candidates []scored

	for _, r := range resources {

		haystack := strings.ToLower(r.Title + " " + r.Content)

		score := 0

		for _, w := range words {
			score += strings.Count(haystack, w)
		}

		if score > 0 {
			candidates = append(candidates, scored{resource: r, score: score})
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score > candidates[j].score
	})

	if len(candidates) > 5 {
		candidates = candidates[:5]
	}

	results := make([]SearchResult, 0, len(candidates))
	var contextBuilder strings.Builder

	for _, cand := range candidates {

		snippet := buildSnippet(cand.resource.Content, words)

		results = append(results, SearchResult{
			ResourceID: cand.resource.ID.String(),
			Title:      cand.resource.Title,
			Snippet:    snippet,
			Score:      cand.score,
		})

		contextBuilder.WriteString("Resource: " + cand.resource.Title + "\n")
		contextBuilder.WriteString(snippet + "\n\n")
	}

	if len(results) == 0 {
		return &SearchResponse{
			Answer:  "I couldn't find anything relevant in your resources yet. Try adding more content first.",
			Results: results,
		}, nil
	}

	answer, err := generateSearchAnswer(contextBuilder.String(), query)

	if err != nil {
		return nil, err
	}

	return &SearchResponse{
		Answer:  answer,
		Results: results,
	}, nil
}

func tokenize(text string) []string {

	raw := strings.Fields(strings.ToLower(text))

	words := make([]string, 0, len(raw))

	for _, w := range raw {

		if stopWords[w] {
			continue
		}

		words = append(words, w)
	}

	return words
}

func buildSnippet(content string, words []string) string {

	lower := strings.ToLower(content)

	idx := -1

	for _, w := range words {
		if i := strings.Index(lower, w); i != -1 {
			idx = i
			break
		}
	}

	if idx == -1 {
		idx = 0
	}

	start := idx - 80
	if start < 0 {
		start = 0
	}

	end := idx + 200
	if end > len(content) {
		end = len(content)
	}

	snippet := strings.TrimSpace(content[start:end])

	return snippet
}

func generateSearchAnswer(contextText string, query string) (string, error) {

	client, err := ai.NewGeminiClient()

	if err != nil {
		return "", err
	}

	prompt := `
You are DevBrain's search assistant.

Answer the question using ONLY the resource excerpts below.
If the excerpts don't contain the answer, say so honestly.
Keep the answer concise (3-5 sentences).

Resource Excerpts:

` + contextText + `

Question:

` + query

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
