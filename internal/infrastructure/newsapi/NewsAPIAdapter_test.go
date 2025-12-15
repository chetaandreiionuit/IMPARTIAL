package newsapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// [RO] Testăm Adaptorul NewsAPI
// Folosim un server HTTP fals pentru a nu apela API-ul real (care costă și e lent).
func TestNewsAPIAdapter_FetchGlobalHeadlines_Success(t *testing.T) {
	// 1. Simulare Server NewsAPI
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificăm dacă cererea e corectă
		assert.Contains(t, r.URL.String(), "/top-headlines")
		assert.Equal(t, "test-api-key", r.URL.Query().Get("apiKey"))

		// Răsfoim un JSON valid de succes
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"status": "ok",
			"totalResults": 2,
			"articles": [
				{"url": "https://cnn.com/article1"},
				{"url": "https://bbc.com/article2"}
			]
		}`))
	}))
	defer mockServer.Close()

	// 2. Configurare Adaptor
	adapter := NewNewsAPIAdapter("test-api-key")
	adapter.baseURL = mockServer.URL // Redirecționăm către serverul fals

	// 3. Execuție
	urls, err := adapter.FetchGlobalHeadlines(context.Background(), "")

	// 4. Verificare
	assert.NoError(t, err)
	assert.Len(t, urls, 2)
	assert.Equal(t, "https://cnn.com/article1", urls[0])
	assert.Equal(t, "https://bbc.com/article2", urls[1])
}

func TestNewsAPIAdapter_FetchGlobalHeadlines_ErrorFromAPI(t *testing.T) {
	// [RO] Testăm cazul când API-ul returnează eroare (ex: 401 Unauthorized)
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"status": "error", "code": "apiKeyMissing"}`))
	}))
	defer mockServer.Close()

	adapter := NewNewsAPIAdapter("bad-key")
	adapter.baseURL = mockServer.URL

	urls, err := adapter.FetchGlobalHeadlines(context.Background(), "")

	assert.Error(t, err)
	assert.Nil(t, urls)
	assert.Contains(t, err.Error(), "NewsAPI a răspuns cu eroare: 401")
}

func TestNewsAPIAdapter_FetchGlobalHeadlines_EmptyKey(t *testing.T) {
	// [RO] Testăm validarea cheii lipsă
	adapter := NewNewsAPIAdapter("")
	urls, err := adapter.FetchGlobalHeadlines(context.Background(), "")

	assert.Error(t, err)
	assert.Nil(t, urls)
	assert.Contains(t, err.Error(), "Lipsă API Key")
}
