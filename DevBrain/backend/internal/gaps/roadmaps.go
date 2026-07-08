package gaps

// Predefined roadmaps per domain. Kept in-memory (zero infra cost) —
// can be moved to a DB table later without changing the API surface.
var Roadmaps = map[string][]string{
	"Go": {
		"variables", "structs", "interfaces", "goroutines", "channels",
		"error handling", "generics", "context",
	},
	"Docker": {
		"containers", "images", "dockerfile", "volumes", "networking",
		"docker compose", "multi-stage builds",
	},
	"Kubernetes": {
		"pods", "deployments", "services", "ingress", "configmaps",
		"secrets", "helm",
	},
	"System Design": {
		"load balancing", "caching", "database sharding", "message queues",
		"microservices", "rate limiting", "cap theorem",
	},
	"Databases": {
		"indexing", "normalization", "transactions", "replication",
		"sharding", "sql", "nosql",
	},
	"JavaScript": {
		"closures", "promises", "async/await", "event loop",
		"prototypes", "hoisting", "modules",
	},
}

func Domains() []string {

	domains := make([]string, 0, len(Roadmaps))

	for name := range Roadmaps {
		domains = append(domains, name)
	}

	return domains
}
