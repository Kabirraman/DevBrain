package chat

import (
	"strings"
	"sort"
	"github.com/Kabirraman/DevBrain/internal/database"
	"github.com/Kabirraman/DevBrain/internal/models"
)
var stopWords = map[string]bool{
	"what": true,
	"is": true,
	"are": true,
	"how": true,
	"does": true,
	"do": true,
	"the": true,
	"a": true,
	"an": true,
	"to": true,
	"of": true,
	"in": true,
	"for": true,
}

func FindRelevantConcepts(
	question string,
) ([]models.Concept, error) {

	var concepts []models.Concept

	if err := database.DB.Find(&concepts).Error; err != nil {
		return nil, err
	}

	rawWords := strings.Fields(
	strings.ToLower(question),
)

var words []string

for _, word := range rawWords {

	if stopWords[word] {
		continue
	}

	words = append(words, word)
}

	type scoredConcept struct {
		Concept models.Concept
		Score   int
	}

	var scored []scoredConcept

	for _, concept := range concepts {

		score := 0

		name := strings.ToLower(concept.Name)

		for _, word := range words {

			if strings.Contains(name, word) {
				score++
			}
		}

		if score > 0 {
			scored = append(scored, scoredConcept{
				Concept: concept,
				Score:   score,
			})
		}
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	var result []models.Concept

	for _, s := range scored {
		result = append(result, s.Concept)
	}

	if len(result) > 5 {
	result = result[:5]
}

return result, nil
}


func FindRelationships(
	concepts []models.Concept,
) ([]models.Relationship, error) {

	var all []models.Relationship

seen := make(map[string]bool)

	for _, concept := range concepts {

		var rels []models.Relationship

		err := database.DB.
			Where(
				"source = ? OR target = ?",
				strings.ToLower(concept.Name),
				strings.ToLower(concept.Name),
			).
			Find(&rels).Error

		if err != nil {
			return nil, err
		}
		for _, rel := range rels {

	key :=
		rel.Source +
			rel.Relation +
			rel.Target

	if seen[key] {
		continue
	}

	seen[key] = true

	all = append(
		all,
		rel,
	)
}
	}

	return all, nil
}

