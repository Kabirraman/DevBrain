package concepts

type RelationshipResponse struct {
	Relationships []RelationshipDTO `json:"relationships"`
}

type RelationshipDTO struct {
	Source   string `json:"source"`
	Relation string `json:"relation"`
	Target   string `json:"target"`
}