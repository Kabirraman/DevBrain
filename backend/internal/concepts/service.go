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
Extract ONLY software engineering and computer science concepts.

Include:
- programming languages
- frameworks
- libraries
- algorithms
- design patterns
- runtime systems
- compilers
- databases
- networking concepts
- cloud technologies
- software architecture concepts
- development tools

Ignore:
- navigation menus
- website links
- cookie notices
- privacy policies
- analytics
- advertisements
- company information
- legal text
- marketing content
- footer content

Return only meaningful technical concepts.
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

func SaveConcepts(concepts []string) error {

	blacklist := map[string]bool{
		"cookies":            true,
		"cookie":             true,
		"privacy":            true,
		"privacy policy":     true,
		"terms of service":   true,
		"traffic":            true,
		"google":             true,
		"google tag manager": true,
		"analytics":          true,
		"advertisement":      true,
		"footer":             true,
	}

	for _, concept := range concepts {

		concept = strings.TrimSpace(strings.ToLower(concept))

		// Skip empty concepts
		if concept == "" {
			continue
		}

		// Skip noisy concepts
		if blacklist[concept] {
			continue
		}

		c := models.Concept{
			Name: concept,
		}

		err := database.DB.
			Where(models.Concept{
				Name: concept,
			}).
			FirstOrCreate(&c).
			Error

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

func ExtractRelationships(
	content string,
) ([]RelationshipDTO, error) {

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
Extract ONLY meaningful technical relationships.

Ignore:
- website navigation
- cookies
- privacy notices
- analytics
- legal text
- advertisements
- marketing content

Good examples:

Go -> USES -> Goroutines

Docker -> RUNS_IN -> Containers

React -> USES -> JSX

Return ONLY valid JSON.
Do NOT use markdown.
Do NOT use triple backticks.

Format:

{
  "relationships":[
    {
      "source":"Go",
      "relation":"USES",
      "target":"Goroutines"
    }
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

	cleaned := cleanJSONResponse(
		resp.Text(),
	)

	var result RelationshipResponse

	err = json.Unmarshal(
		[]byte(cleaned),
		&result,
	)

	if err != nil {
		return nil, err
	}

	return result.Relationships, nil
}
func SaveRelationships(
	relationships []RelationshipDTO,
) error {

	for _, rel := range relationships {

		source := strings.TrimSpace(
			strings.ToLower(rel.Source),
		)

		relation := strings.TrimSpace(
			strings.ToUpper(rel.Relation),
		)

		target := strings.TrimSpace(
			strings.ToLower(rel.Target),
		)

		r := models.Relationship{
			Source:   source,
			Relation: relation,
			Target:   target,
		}

		err := database.DB.
			Where(models.Relationship{
				Source:   source,
				Relation: relation,
				Target:   target,
			}).
			FirstOrCreate(&r).
			Error

		if err != nil {
			return err
		}
	}

	return nil
}
