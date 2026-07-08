package concepts

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kabirraman/DevBrain/internal/database"
	"github.com/Kabirraman/DevBrain/internal/models"
	"github.com/google/uuid"
	"google.golang.org/genai"
	"os"
	"strings"
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
		fmt.Println("GEMINI API ERROR (ExtractConcepts):", err)
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

func SaveUserConcepts(userID string, concepts []string) error {

	uid, err := uuid.Parse(userID)

	if err != nil {
		return err
	}

	for _, concept := range concepts {

		name := strings.TrimSpace(strings.ToLower(concept))

		if name == "" {
			continue
		}

		var c models.Concept

		err := database.DB.
			Where(models.Concept{Name: name}).
			FirstOrCreate(&c).
			Error

		if err != nil {
			return err
		}

		uc := models.UserConcept{
			UserID:    uid,
			ConceptID: c.ID,
		}

		err = database.DB.
			Where(models.UserConcept{
				UserID:    uid,
				ConceptID: c.ID,
			}).
			FirstOrCreate(&uc).
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

	text = strings.TrimSpace(text)

	// Gemini sometimes wraps the JSON in explanatory text even when told
	// not to. Extract just the outermost {...} block as a safety net.
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")

	if start != -1 && end != -1 && end > start {
		text = text[start : end+1]
	}

	return text
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
You are building a knowledge graph.

Extract important concepts and relationships from the text.
The goal is to create a connected knowledge graph.
Prefer relationships that connect concepts together.
Avoid isolated concepts whenever possible.

Rules:

1. Return ONLY valid JSON.
2. Create relationships between major concepts whenever hierarchy, ownership, versioning, dependency, inclusion, or usage is implied.
3. Prefer connecting concepts instead of leaving them isolated.
4. Use concise relationship names in UPPERCASE.

Examples:

Go HAS_RELEASE Go 1.26
Go PROVIDES pkg.go.dev API
Go HAS_TOOL Playground
Go HAS_TOOL go fix
Go INCLUDES Standard Library

Return ONLY valid JSON.
Only extract relationships that are useful for learning,
understanding, programming, software architecture,
documentation, tooling, APIs, versions, concepts,
frameworks, algorithms, or technical knowledge.

Ignore:

- Social media accounts
- Community links
- Cookies
- Navigation menus
- Help pages
- Download buttons
- Contact links
- Footer content
- Branding links
- Legal pages

Output format:

{
  "relationships": [
    {
      "source": "Go",
      "relation": "HAS_RELEASE",
      "target": "Go 1.26"
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
		fmt.Println("GEMINI API ERROR (ExtractRelationships):", err)
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
	userID string,
	relationships []RelationshipDTO,
) error {

	uid, err := uuid.Parse(userID)

	if err != nil {
		return err
	}

	var blacklist = map[string]bool{
		"twitter":        true,
		"slack":          true,
		"mastodon":       true,
		"bluesky":        true,
		"help":           true,
		"download":       true,
		"stack overflow": true,
	}

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

		// Skip empty values
		if source == "" ||
			relation == "" ||
			target == "" {
			continue
		}

		// Skip self-loops
		if source == target {
			continue
		}

		// Skip noisy concepts
		if blacklist[source] ||
			blacklist[target] {
			continue
		}

		r := models.Relationship{
			UserID:   uid,
			Source:   source,
			Relation: relation,
			Target:   target,
		}

		err := database.DB.
			Where(models.Relationship{
				UserID:   uid,
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