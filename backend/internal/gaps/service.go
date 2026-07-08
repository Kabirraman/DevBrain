package gaps

import (
	"errors"
	"strings"

	"github.com/Kabirraman/DevBrain/internal/database"
	"github.com/Kabirraman/DevBrain/internal/models"
)

func GetGapAnalysis(userID string, domain string) (*GapResponse, error) {

	roadmap, ok := Roadmaps[domain]

	if !ok {
		return nil, errors.New("unknown domain")
	}

	knownSet, err := GetKnownConceptSet(userID)

	if err != nil {
		return nil, err
	}

	return ComputeGap(domain, roadmap, knownSet), nil
}

// ComputeGap is a pure function so callers that need multiple domains
// (e.g. the analytics dashboard) can fetch the known-concept set once and
// reuse it, instead of re-querying the database for every domain.
func ComputeGap(domain string, roadmap []string, knownSet map[string]bool) *GapResponse {

	known := []string{}
	missing := []string{}

	for _, topic := range roadmap {

		if knownSet[strings.ToLower(topic)] {
			known = append(known, topic)
		} else {
			missing = append(missing, topic)
		}
	}

	completion := 0

	if len(roadmap) > 0 {
		completion = (len(known) * 100) / len(roadmap)
	}

	return &GapResponse{
		Domain:     domain,
		Completion: completion,
		Known:      known,
		Missing:    missing,
	}
}

func GetKnownConceptSet(userID string) (map[string]bool, error) {

	var userConcepts []models.UserConcept

	err := database.DB.
		Where("user_id = ?", userID).
		Find(&userConcepts).
		Error

	if err != nil {
		return nil, err
	}

	set := make(map[string]bool)

	if len(userConcepts) == 0 {
		return set, nil
	}

	conceptIDs := make([]string, 0, len(userConcepts))

	for _, uc := range userConcepts {
		conceptIDs = append(conceptIDs, uc.ConceptID.String())
	}

	var concepts []models.Concept

	err = database.DB.
		Where("id IN ?", conceptIDs).
		Find(&concepts).
		Error

	if err != nil {
		return nil, err
	}

	for _, c := range concepts {
		set[strings.ToLower(c.Name)] = true
	}

	return set, nil
}
