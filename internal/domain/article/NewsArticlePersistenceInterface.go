package article

import (
	"context"

	"github.com/google/uuid"
)

// [RO] Interfața de Persistență a Știrilor
//
// Această interfață definește "Contractul" pe care trebuie să îl respecte orice
// sistem de stocare (Bază de Date) pe care îl folosim.
//
// Ea spune CE trebuie să facă baza de date, nu CUM să o facă.
// (Ex: Trebuie să poți salva o știre, nu contează dacă o scrii într-un fișier sau în Postgres).
type NewsArticlePersistenceInterface interface {
	// [RO] Salvează Știrea
	// Scrie permanent articolul și toate relațiile lui în baza de date.
	// Returnează o eroare dacă operațiunea eșuează.
	PersistNewsArticle(execution_context context.Context, article *NewsArticleEntity) error

	// [RO] Găsește Știrea după ID
	// Recuperează dosarul complet al știrii folosind codul său unic (UUID).
	RetrieveNewsArticleByID(execution_context context.Context, id uuid.UUID) (*NewsArticleEntity, error)

	// [RO] Caută Știri Similare (Semantic)
	// Folosește vectori matematici pentru a găsi alte articole care vorbesc despre
	// același subiect, chiar dacă folosesc cuvinte diferite.
	// (Ex: "Război" ~ "Conflict armat").
	FindSemanticallySimilarArticles(execution_context context.Context, embedding []float32, limit int) ([]*NewsArticleEntity, error)

	// [RO] Verifică Existența (Deduplicare)
	// O verificare rapidă pentru a vedea dacă acest URL a mai fost procesat vreodată.
	// Folosită pentru a nu consuma credite AI pe același articol de două ori.
	CheckIfArticleExistsByURL(execution_context context.Context, url string) (bool, error)
}
