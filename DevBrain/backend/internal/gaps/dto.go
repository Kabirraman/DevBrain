package gaps

type GapResponse struct {
	Domain     string   `json:"domain"`
	Completion int      `json:"completion"`
	Known      []string `json:"known"`
	Missing    []string `json:"missing"`
}

type DomainsResponse struct {
	Domains []string `json:"domains"`
}
