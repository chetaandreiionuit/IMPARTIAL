package scraper

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-shiori/go-readability"
	"github.com/gocolly/colly/v2"
)

// [RO] Date Extrase
type ScrapedData struct {
	Title       string
	CleanText   string
	Excerpt     string
	Author      string
	PublishDate string
}

// [RO] Serviciul de Scraping
type CollyScraper struct {
	collector *colly.Collector
}

func NewCollyScraper() *CollyScraper {
	// [RO] Configurare Colly
	// 1. Singleton Collector pentru eficiență
	c := colly.NewCollector(
		colly.Async(true), // Execuție asincronă (Fan-out)
		colly.UserAgent("TruthWeaveBot/2.0 (+http://truthweave.internal/bot)"),
	)

	// [RO] Limitare Rată & Politețe (Architect-Alpha Spec)
	// - DomainGlob: * (pentru toate domeniile)
	// - Parallelism: 2 (maxim 2 request-uri simultane per domeniu)
	// - RandomDelay: 5s (jitter pentru a evita detecția)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	return &CollyScraper{
		collector: c,
	}
}

// [RO] Scrape URL
// Această funcție este blocantă per-URL, dar intern Colly poate paralela.
// Deoarece folosim `readability` pe `ctx.Body`, trebuie să fim atenți la callback.
func (s *CollyScraper) Scrape(targetURL string) (*ScrapedData, error) {
	var result *ScrapedData
	var errResult error

	// Clona colectorului pentru request individual izolat (preferabil pentru control context)
	// Dar LimitRule e pe instanța principală. Folosim instanța principală.
	c := s.collector.Clone()

	// [RO] Callback - La primirea răspunsului
	c.OnResponse(func(r *colly.Response) {
		// [RO] Curățare Semantică DIRECT din memorie (fără re-descărcare)
		// r.Body este []byte. Il convertim in Reader.
		article, err := readability.FromReader(strings.NewReader(string(r.Body)), r.Request.URL)
		if err != nil {
			errResult = fmt.Errorf("readability failed: %w", err)
			return
		}

		result = &ScrapedData{
			Title:       article.Title,
			CleanText:   article.TextContent,
			Excerpt:     article.Excerpt,
			Author:      article.Byline,
			PublishDate: "", // Readability nu extrage data meta fiabil mereu
		}
	})

	c.OnError(func(_ *colly.Response, err error) {
		errResult = err
	})

	// Blocăm până termină (pentru acest URL simplu)
	// Deoarece avem Async(true), .Visit pornește goroutine. .Wait() așteaptă.
	// Dar .Clone() moștenește setările? Da.
	// Facem un truc: pentru un singur URL sincron (chemat din Temporal Activity),
	// nu neaparat avem nevoie de Async pe collector, ci Temporal face async.
	// Dar LimitRule cere Async in Colly v2 pentru a funcționa corect pe mai multe threaduri?
	// LimitRule blochează thread-ul apelant dacă e Sync.

	// Pentru simplitate și integrare cu Temporal, aici executăm sincron vizita.
	// Temporal se ocupă de paralelism (lansaând mai multe Activități).
	// LimitRule din Colly va bloca activitatea Temporal dacă depășim rata.

	err := c.Visit(targetURL)
	if err != nil {
		return nil, err
	}

	c.Wait() // Așteptăm să termine procesarea (OnResponse)

	if errResult != nil {
		return nil, errResult
	}
	if result == nil {
		return nil, fmt.Errorf("no content extracted")
	}

	return result, nil
}
