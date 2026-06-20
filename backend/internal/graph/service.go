package graph

import (
	"github.com/Kabirraman/DevBrain/internal/database"
	"github.com/Kabirraman/DevBrain/internal/models"
)

func BuildGraph() (*GraphResponse, error) {

	var concepts []models.Concept
	var relationships []models.Relationship

	if err := database.DB.Find(&concepts).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Find(&relationships).Error; err != nil {
		return nil, err
	}

	nodes := []GraphNode{}
	edges := []GraphEdge{}

	// Prevent duplicate nodes
	nodeMap := make(map[string]bool)

	// Add concept nodes
	for _, concept := range concepts {

		if !nodeMap[concept.Name] {

			nodes = append(nodes, GraphNode{
				ID: concept.Name,
			})

			nodeMap[concept.Name] = true
		}
	}

	// Prevent duplicate edges
	edgeMap := make(map[string]bool)

	for _, rel := range relationships {

		// Ensure source node exists
		if !nodeMap[rel.Source] {

			nodes = append(nodes, GraphNode{
				ID: rel.Source,
			})

			nodeMap[rel.Source] = true
		}

		// Ensure target node exists
		if !nodeMap[rel.Target] {

			nodes = append(nodes, GraphNode{
				ID: rel.Target,
			})

			nodeMap[rel.Target] = true
		}

		edgeKey := rel.Source + "|" + rel.Relation + "|" + rel.Target

		if !edgeMap[edgeKey] {

			edges = append(edges, GraphEdge{
				Source: rel.Source,
				Target: rel.Target,
				Label:  rel.Relation,
			})

			edgeMap[edgeKey] = true
		}
	}

	return &GraphResponse{
		Nodes: nodes,
		Edges: edges,
	}, nil
}