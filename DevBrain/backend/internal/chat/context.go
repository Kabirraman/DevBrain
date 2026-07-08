package chat

import (
	"strings"

	"github.com/Kabirraman/DevBrain/internal/models"
)
func BuildContext(
	concepts []models.Concept,
	relationships []models.Relationship,
) string {

	var builder strings.Builder

	builder.WriteString(
		"Relevant Concepts:\n\n",
	)

	for _, c := range concepts {

		builder.WriteString(
			"- " + c.Name + "\n",
		)
	}

	builder.WriteString(
		"\nRelationships:\n\n",
	)

	for _, r := range relationships {

		builder.WriteString(
			r.Source +
				" " +
				r.Relation +
				" " +
				r.Target +
				"\n",
		)
	}

	return builder.String()
}