package analytics

import (
	"sort"

	"github.com/Kabirraman/DevBrain/internal/database"
	"github.com/Kabirraman/DevBrain/internal/gaps"
	"github.com/Kabirraman/DevBrain/internal/models"
)

func GetAnalytics(userID string) (*AnalyticsResponse, error) {

	var totalResources int64
	var userConcepts []models.UserConcept

	if err := database.DB.
		Model(&models.Resource{}).
		Where("user_id = ?", userID).
		Count(&totalResources).Error; err != nil {
		return nil, err
	}

	if err := database.DB.
		Where("user_id = ?", userID).
		Find(&userConcepts).Error; err != nil {
		return nil, err
	}

	conceptIDs := make([]string, 0, len(userConcepts))

	for _, uc := range userConcepts {
		conceptIDs = append(conceptIDs, uc.ConceptID.String())
	}

	var concepts []models.Concept

	if len(conceptIDs) > 0 {
		if err := database.DB.
			Where("id IN ?", conceptIDs).
			Find(&concepts).Error; err != nil {
			return nil, err
		}
	}

	names := make([]string, 0, len(concepts))
	for _, c := range concepts {
		names = append(names, c.Name)
	}

	sort.Strings(names)

	topTopics := names
	if len(topTopics) > 5 {
		topTopics = topTopics[:5]
	}

	coverage := []CoverageEntry{}
	missing := map[string]bool{}

	for domain := range gaps.Roadmaps {

		analysis, err := gaps.GetGapAnalysis(userID, domain)

		if err != nil {
			continue
		}

		coverage = append(coverage, CoverageEntry{
			Domain:     domain,
			Completion: analysis.Completion,
		})

		for _, m := range analysis.Missing {
			missing[m] = true
		}
	}

	sort.Slice(coverage, func(i, j int) bool {
		return coverage[i].Domain < coverage[j].Domain
	})

	missingList := make([]string, 0, len(missing))
	for m := range missing {
		missingList = append(missingList, m)
	}
	sort.Strings(missingList)

	if len(missingList) > 8 {
		missingList = missingList[:8]
	}

	return &AnalyticsResponse{
		TotalResources:    int(totalResources),
		TotalConcepts:     len(concepts),
		MostLearnedTopics: topTopics,
		KnowledgeCoverage: coverage,
		MissingConcepts:   missingList,
	}, nil
}
