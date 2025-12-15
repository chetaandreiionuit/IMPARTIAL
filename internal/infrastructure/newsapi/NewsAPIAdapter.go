package newsapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// [RO] Adaptor pentru NewsAPI.org
// Această componentă este "Corespondentul Străin" care ne aduce știrile din lume.
type NewsAPIAdapter struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// [RO] Constructor NewsAPI
func NewNewsAPIAdapter(apiKey string) *NewsAPIAdapter {
	return &NewsAPIAdapter{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    "https://newsapi.org/v2",
	}
}

// Structuri interne pentru răspunsul JSON de la NewsAPI
type newsAPIResponse struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		URL string `json:"url"`
	} `json:"articles"`
}

// [RO] Aduce Ultimele Titluri (FetchGlobalHeadlines)
// Caută știri globale (top-headlines) sau pe o temă anume.
// Returnează o listă de URL-uri brute.
func (adapter *NewsAPIAdapter) FetchGlobalHeadlines(ctx context.Context, query string) ([]string, error) {
	// Dacă nu avem cheie, returnăm o listă goală (sau eroare, dar pentru reziliență, doar logăm).
	if adapter.apiKey == "" {
		return nil, fmt.Errorf("[RO] Lipsă API Key pentru NewsAPI. Nu pot aduce știri.")
	}

	// Construim URL-ul cererii
	// Implicit căutăm știri generale în limba engleză (pentru acoperire globală maximă).
	// Dacă 'query' e gol, luăm 'general'.
	endpoint := fmt.Sprintf("%s/top-headlines?language=en&pageSize=20&apiKey=%s", adapter.baseURL, adapter.apiKey)
	if query != "" {
		endpoint += "&q=" + query
	}

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creare cerere eșuată: %w", err)
	}

	resp, err := adapter.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("apel rețea eșuat: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NewsAPI a răspuns cu eroare: %d", resp.StatusCode)
	}

	var parsedResp newsAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsedResp); err != nil {
		return nil, fmt.Errorf("decodare JSON eșuată: %w", err)
	}

	var urls []string
	for _, art := range parsedResp.Articles {
		if art.URL != "" {
			urls = append(urls, art.URL)
		}
	}

	return urls, nil
}
