package views

type OmdbSearchResult struct {
	Search       []OmdbSearchMovie `json:"Search"`
	TotalResults string            `json:"totalResults"`
	Response     string            `json:"Response"`
}
