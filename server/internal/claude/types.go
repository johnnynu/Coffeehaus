package claude

// SearchIntent represents the structured output from Claude
type SearchIntent struct {
    SearchType      string       `json:"searchType"`
    NormalizedQuery string       `json:"normalizedQuery"`
    Location        *Location    `json:"location,omitempty"`
    Terms          SearchTerms   `json:"terms"`
}

type Location struct {
    Name   string  `json:"name,omitempty"`
    Radius float64 `json:"radius,omitempty"`
}

type SearchTerms struct {
    Shop    string   `json:"shop,omitempty"`
    Filters []string `json:"filters"`
}