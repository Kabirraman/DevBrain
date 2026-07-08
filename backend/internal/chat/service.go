package chat
import (
	"github.com/Kabirraman/DevBrain/internal/models"
	
)

func AskGraph(
	userID string,
	question string,
) (string, []string, error) {

	concepts, err := FindRelevantConcepts(userID, question)

if err != nil {
	return "", nil, err
}

relationships, err := FindRelationships(userID, concepts)

if err != nil {
	return "", nil, err
}

context := BuildContext(
	concepts,
	relationships,
)

answer, err := GenerateAnswer(
	context,
	question,
)

if err != nil {
	return "", nil, err
}

return answer,
	extractConceptNames(concepts),
	nil
}

func extractConceptNames(
	concepts []models.Concept,
) []string {

	var names []string

	for _, c := range concepts {

		names = append(
			names,
			c.Name,
		)
	}

	return names
}