package gdelt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// [RO] Constante GDELT
const (
	GdeltApiBaseURL = "https://api.gdeltproject.org/api/v2/doc/doc"
)

// [RO] Structura Articol GDELT (conform specificației API DOC 2.0)
type GdeltArticle struct {
	URL           string `json:"url"`
	Title         string `json:"title"`
	Seendate      string `json:"seendate"` // Format: YYYYMMDDHHMMSS
	SourceCountry string `json:"sourcecountry"`
	Language      string `json:"language"`
	SocialImage   string `json:"socialimage"`
	Domain        string `json:"domain"`
}

// [RO] Wrapper Răspuns GDELT
type GdeltApiResponse struct {
	Articles []GdeltArticle `json:"articles"`
}

// [RO] Filtre de Interogare
type GdeltQueryFilters struct {
	ToneAbsGreaterThan float64 // toneabs>5
	ImageTag           string  // imagetag:"tank"
	SourceLang         string  // sourcelang:rum
	Theme              string  // theme:CYBER_ATTACK
	Keywords           string  // "NATO Russia"
}

// [RO] Adaptor GDELT
type GDELTAdapter struct {
	client *http.Client
}

func NewGDELTAdapter() *GDELTAdapter {
	return &GDELTAdapter{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// [RO] Construiește URL-ul de interogare
func (a *GDELTAdapter) BuildGDELTQuery(filters GdeltQueryFilters) string {
	// Query de bază
	queryParts := ""

	if filters.Keywords != "" {
		queryParts += filters.Keywords
	}

	if filters.ToneAbsGreaterThan > 0 {
		queryParts += fmt.Sprintf(" toneabs>%2.f", filters.ToneAbsGreaterThan)
	}

	if filters.ImageTag != "" {
		queryParts += fmt.Sprintf(" imagetag:%q", filters.ImageTag)
	}

	if filters.SourceLang != "" {
		queryParts += fmt.Sprintf(" sourcelang:%s", filters.SourceLang)
	}

	if filters.Theme != "" {
		queryParts += fmt.Sprintf(" theme:%s", filters.Theme)
	}

	// Parametrii standard: JSON, ArtList, MaxRecords, Timespan
	// Timespan 15min pentru a prinde fluxul curent
	params := url.Values{}
	params.Add("query", queryParts)
	params.Add("mode", "artlist")
	params.Add("format", "json")
	params.Add("maxrecords", "50")  // Limităm la 50 per batch
	params.Add("timespan", "15min") // Doar ultimele 15 minute

	return fmt.Sprintf("%s?%s", GdeltApiBaseURL, params.Encode())
}

// [RO] Fetch Articles
func (a *GDELTAdapter) FetchLatestArticles(ctx context.Context, filters GdeltQueryFilters) ([]GdeltArticle, error) {
	fullURL := a.BuildGDELTQuery(filters)

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GDELT API error: %d", resp.StatusCode)
	}

	var result GdeltApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		// Uneori GDELT returneaza erori ciudate daca nu sunt rezultate, dar format=json ar trebui sa fie consistent.
		return nil, err
	}

	return result.Articles, nil
}
