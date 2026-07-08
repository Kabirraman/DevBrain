package concepts

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kabirraman/DevBrain/internal/database"
	"github.com/Kabirraman/DevBrain/internal/models"
	"github.com/google/uuid"
	"google.golang.org/genai"
	"gorm.io/gorm/clause"
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

var conceptBlacklist = map[string]bool{
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

// normalizeConceptNames lowercases, trims, dedupes, and filters out noise -
// shared by SaveConcepts and SaveUserConcepts so both agree on what counts
// as a valid concept name.
func normalizeConceptNames(concepts []string) []string {

	seen := make(map[string]bool)
	names := make([]string, 0, len(concepts))

	for _, concept := range concepts {

		name := strings.TrimSpace(strings.ToLower(concept))

		if name == "" || conceptBlacklist[name] || seen[name] {
			continue
		}

		seen[name] = true
		names = append(names, name)
	}

	return names
}

// upsertConcepts ensures every given name exists in the concepts table using
// at most 2 queries total (one lookup + one batch insert), instead of one
// round trip per concept - this matters a lot when a single page extracts
// hundreds of concepts against a remote database. Returns name -> ID.
func upsertConcepts(names []string) (map[string]uuid.UUID, error) {

	result := make(map[string]uuid.UUID)

	if len(names) == 0 {
		return result, nil
	}

	var existing []models.Concept

	if err := database.DB.
		Where("name IN ?", names).
		Find(&existing).
		Error; err != nil {
		return nil, err
	}

	for _, c := range existing {
		result[c.Name] = c.ID
	}

	var toCreate []models.Concept

	for _, name := range names {
		if _, ok := result[name]; !ok {
			toCreate = append(toCreate, models.Concept{Name: name})
		}
	}

	if len(toCreate) > 0 {

		if err := database.DB.
			Clauses(clause.OnConflict{DoNothing: true}).
			Create(&toCreate).
			Error; err != nil {
			return nil, err
		}

		// Rows that hit a conflict (already existed) keep the bogus
		// in-memory ID assigned by BeforeCreate before the insert was
		// skipped - re-fetch to get the authoritative IDs for all names.
		var confirmed []models.Concept

		if err := database.DB.
			Where("name IN ?", names).
			Find(&confirmed).
			Error; err != nil {
			return nil, err
		}

		for _, c := range confirmed {
			result[c.Name] = c.ID
		}
	}

	return result, nil
}

func SaveConcepts(concepts []string) error {

	names := normalizeConceptNames(concepts)

	_, err := upsertConcepts(names)

	return err
}

func SaveUserConcepts(userID string, concepts []string) error {

	uid, err := uuid.Parse(userID)

	if err != nil {
		return err
	}

	names := normalizeConceptNames(concepts)

	conceptIDs, err := upsertConcepts(names)

	if err != nil {
		return err
	}

	ids := make([]uuid.UUID, 0, len(conceptIDs))

	for _, id := range conceptIDs {
		ids = append(ids, id)
	}

	var existingLinks []models.UserConcept

	if len(ids) > 0 {
		if err := database.DB.
			Where("user_id = ? AND concept_id IN ?", uid, ids).
			Find(&existingLinks).
			Error; err != nil {
			return err
		}
	}

	linked := make(map[uuid.UUID]bool)
	for _, l := range existingLinks {
		linked[l.ConceptID] = true
	}

	var toCreate []models.UserConcept

	for _, id := range ids {
		if !linked[id] {
			toCreate = append(toCreate, models.UserConcept{
				UserID:    uid,
				ConceptID: id,
			})
		}
	}

	if len(toCreate) > 0 {
		if err := database.DB.
			Clauses(clause.OnConflict{DoNothing: true}).
			Create(&toCreate).
			Error; err != nil {
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

	type triple struct {
		source   string
		relation string
		target   string
	}

	seen := make(map[triple]bool)
	var clean []triple

	for _, rel := range relationships {

		source := strings.TrimSpace(strings.ToLower(rel.Source))
		relation := strings.TrimSpace(strings.ToUpper(rel.Relation))
		target := strings.TrimSpace(strings.ToLower(rel.Target))

		if source == "" || relation == "" || target == "" {
			continue
		}

		if source == target {
			continue
		}

		if blacklist[source] || blacklist[target] {
			continue
		}

		t := triple{source, relation, target}

		if seen[t] {
			continue
		}

		seen[t] = true
		clean = append(clean, t)
	}

	if len(clean) == 0 {
		return nil
	}

	var existing []models.Relationship

	if err := database.DB.
		Where("user_id = ?", uid).
		Find(&existing).
		Error; err != nil {
		return err
	}

	existingSet := make(map[triple]bool, len(existing))

	for _, e := range existing {
		existingSet[triple{e.Source, e.Relation, e.Target}] = true
	}

	var toCreate []models.Relationship

	for _, t := range clean {

		if existingSet[t] {
			continue
		}

		toCreate = append(toCreate, models.Relationship{
			UserID:   uid,
			Source:   t.source,
			Relation: t.relation,
			Target:   t.target,
		})
	}

	if len(toCreate) == 0 {
		return nil
	}

	return database.DB.
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&toCreate).
		Error
}