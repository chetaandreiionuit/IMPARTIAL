package article

import (
	"testing"
)

// [RO] Teste Unitare pentru Entitatea Știre
func TestNewsArticleEntity_VerifyDataIntegrity(t *testing.T) {
	tests := []struct {
		name    string
		article NewsArticleEntity
		wantErr bool
	}{
		{
			name: "[RO] Caz Valid: Articol cu URL Corect",
			article: NewsArticleEntity{
				OriginalURL: "https://example.com/news",
				Content:     "Conținut valid",
			},
			wantErr: false,
		},
		{
			name: "[RO] Caz Invalid: URL Lipsă",
			article: NewsArticleEntity{
				OriginalURL: "",
				Content:     "Conținut orfan",
			},
			wantErr: true,
		},
		{
			name: "[RO] Caz Invalid: Format URL Greșit",
			article: NewsArticleEntity{
				OriginalURL: "nu-este-un-link-valid",
				Content:     "Text oarecare",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apelăm metoda redenumită 'VerifyDataIntegrity'
			if err := tt.article.VerifyDataIntegrity(); (err != nil) != tt.wantErr {
				t.Errorf("VerifyDataIntegrity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
