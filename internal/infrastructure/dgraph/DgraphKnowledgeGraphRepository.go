package dgraph

import (
	"context"
	"encoding/json"
	"fmt" // "fmt" was used in original
	"time"

	"github.com/dgraph-io/dgo/v240"
	"github.com/dgraph-io/dgo/v240/protos/api"
	"github.com/yourorg/truthweave/internal/domain/article"
)

// [RO] Depozit Graf de Cunoștințe (Dgraph)
//
// Această componentă gestionează "Creierul Asociativ" al aplicației.
// Spre deosebire de Postgres (care e ca un tabel Excel), Dgraph funcționează ca o rețea neuronală,
// legând concepte între ele (ex: "Articolul A" -> menționează "Persoana X" -> care apare și în "Articolul B").
type DgraphKnowledgeGraphRepository struct {
	graphClient *dgo.Dgraph
}

// [RO] Constructor Graf
func NewDgraphKnowledgeGraphRepository(client *dgo.Dgraph) *DgraphKnowledgeGraphRepository {
	return &DgraphKnowledgeGraphRepository{graphClient: client}
}

// [RO] Salvează Știrea în Graf (Implementare)
//
// Această metodă nu doar salvează textul, ci creează "noduri" și "muchii" în graf.
// Dacă articolul menționează "București", sistemul va crea (sau refolosi) nodul "București"
// și va trage o linie între Articol și Oraș.
func (repo *DgraphKnowledgeGraphRepository) SaveNewsArticleToGraph(executionContext context.Context, newsArticle *article.NewsArticleEntity) error {
	transaction := repo.graphClient.NewTxn()
	defer transaction.Discard(executionContext)

	// [RO] Pasul 1: Gestionarea Entităților Menționate
	// Înainte să salvăm articolul, ne asigurăm că toate persoanele/locurile menționate
	// au un ID în graf. Dacă nu există, le creăm.
	entityUIDs := make([]map[string]string, 0)

	for _, ent := range newsArticle.Mentions {
		// Încercăm să găsim entitatea existentă
		query := `query q($name: string) {
			ents(func: eq(name, $name)) {
				uid
			}
		}`

		response, err := transaction.QueryWithVars(executionContext, query, map[string]string{"$name": ent.Name})
		if err != nil {
			return fmt.Errorf("failed to query entity %s: %w", ent.Name, err)
		}

		var root struct {
			Ents []struct {
				Uid string `json:"uid"`
			} `json:"ents"`
		}
		if err := json.Unmarshal(response.Json, &root); err != nil {
			return err
		}

		var uid string
		if len(root.Ents) > 0 {
			uid = root.Ents[0].Uid
		} else {
			// [RO] Creare Nod Nou
			// Dacă entitatea nu există, o creăm acum.
			nquads := fmt.Sprintf(`_:new <name> "%s" . _:new <type> "%s" .`, ent.Name, ent.Type)
			mutation := &api.Mutation{SetNquads: []byte(nquads)}
			assigned, err := transaction.Mutate(executionContext, mutation)
			if err != nil {
				return err
			}
			uid = assigned.Uids["new"]
		}

		entityUIDs = append(entityUIDs, map[string]string{"uid": uid})
	}

	// [RO] Pasul 2: Salvarea Articolului și a Legăturilor

	// Definim structura pentru serializare JSON (DTO) compatibilă cu Dgraph
	type ArticleGraphDTO struct {
		Uid         string              `json:"uid,omitempty"`
		URL         string              `json:"url"`
		Title       string              `json:"title"`
		Score       float64             `json:"truth_score"`
		Mentions    []map[string]string `json:"mentioned_entities"`
		ArweaveID   string              `json:"arweave_tx_id"`
		SolanaSig   string              `json:"solana_signature"`
		PublishedAt string              `json:"published_at"`
		DType       []string            `json:"dgraph.type,omitempty"`
	}

	// Căutăm dacă articolul există deja în graf (pentru update)
	qArticle := `query q($url: string) {
		arts(func: eq(url, $url)) {
			uid
		}
	}`
	res, err := transaction.QueryWithVars(executionContext, qArticle, map[string]string{"$url": newsArticle.OriginalURL})
	if err != nil {
		return err
	}

	var artRoot struct {
		Arts []struct {
			Uid string `json:"uid"`
		} `json:"arts"`
	}
	json.Unmarshal(res.Json, &artRoot)

	var articleUID string
	if len(artRoot.Arts) > 0 {
		articleUID = artRoot.Arts[0].Uid
	}

	dto := ArticleGraphDTO{
		Uid:         articleUID,
		URL:         newsArticle.OriginalURL,
		Title:       newsArticle.Title,
		Score:       newsArticle.TruthScore,
		Mentions:    entityUIDs, // Aici facem legătura fizică în graf!
		ArweaveID:   newsArticle.ArweaveTransactionID,
		SolanaSig:   newsArticle.SolanaSignature,
		PublishedAt: newsArticle.PublishedAt.Format(time.RFC3339),
		DType:       []string{"Article"},
	}

	jsonData, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	mutation := &api.Mutation{
		SetJson:   jsonData,
		CommitNow: true, // Salvăm totul atomic (totul sau nimic)
	}

	_, err = transaction.Mutate(executionContext, mutation)
	if err != nil {
		return fmt.Errorf("failed to save article graph: %w", err)
	}

	return nil
}

// [RO] Verifică Existența (Pentru Graf)
// Deși verificăm în Postgres, uneori e util să verificăm direct în graf.
func (repo *DgraphKnowledgeGraphRepository) CheckIfArticleExistsInGraph(executionContext context.Context, url string) (bool, error) {
	transaction := repo.graphClient.NewReadOnlyTxn()
	const query = `query val($url: string) {
		cnt(func: eq(url, $url)) {
			count(uid)
		}
	}`

	resp, err := transaction.QueryWithVars(executionContext, query, map[string]string{"$url": url})
	if err != nil {
		return false, err
	}

	var root struct {
		Cnt []struct {
			Count int `json:"count"`
		} `json:"cnt"`
	}
	if err := json.Unmarshal(resp.Json, &root); err != nil {
		return false, err
	}
	return len(root.Cnt) > 0 && root.Cnt[0].Count > 0, nil
}

// [RO] Creează Muchie Cauzală
// Stabilește o legătură de tip "caused_by" între două evenimente (articole).
func (repo *DgraphKnowledgeGraphRepository) CreateCausalEdge(executionContext context.Context, parentID string, childID string, relationType string) error {
	// Presupunem că ID-urile sunt UUID-uri pe care le mapăm la UIDs interne.
	// În Dgraph schema avem 'event.id' exact index.

	// 1. Găsim UIDs interne pentru Parent și Child
	query := `query q($pid: string, $cid: string) {
		parent(func: eq(event.id, $pid)) {
			uid
		}
		child(func: eq(event.id, $cid)) {
			uid
		}
	}`

	transaction := repo.graphClient.NewTxn()
	// Nota: Nu facem defer Discard imediat pentru că vrem CommitNow sau manual commit.
	// Dar e safe să facem defer Discard, commit va face discard no-op.
	defer transaction.Discard(executionContext)

	resp, err := transaction.QueryWithVars(executionContext, query, map[string]string{
		"$pid": parentID,
		"$cid": childID,
	})
	if err != nil {
		return fmt.Errorf("failed to query UIDs for causal edge: %w", err)
	}

	var root struct {
		Parent []struct {
			Uid string `json:"uid"`
		} `json:"parent"`
		Child []struct {
			Uid string `json:"uid"`
		} `json:"child"`
	}
	if err := json.Unmarshal(resp.Json, &root); err != nil {
		return err
	}

	if len(root.Parent) == 0 || len(root.Child) == 0 {
		return fmt.Errorf("parent or child event not found in graph")
	}

	parentUid := root.Parent[0].Uid
	childUid := root.Child[0].Uid

	// 2. Creăm muchia (Child --[caused_by]--> Parent)
	// Adăugăm și fațete pe muchie (relația, greutatea) dacă e suportat,
	// momentan simplu link.
	nquad := fmt.Sprintf(`<%s> <event.caused_by> <%s> .`, childUid, parentUid)

	mutation := &api.Mutation{
		SetNquads: []byte(nquad),
		CommitNow: true,
	}

	if _, err := transaction.Mutate(executionContext, mutation); err != nil {
		return fmt.Errorf("failed to create causal edge: %w", err)
	}

	return nil
}

// [RO] Upsert Causal Event (Part 2 Refactoring)
// Salvează rezultatul analizei cauzale (nodul + scorurile).
func (repo *DgraphKnowledgeGraphRepository) UpsertCausalEvent(ctx context.Context, eventID string, timestamp time.Time, summary string, score float64) error {
	type EventDTO struct {
		Uid        string   `json:"uid"`
		DType      []string `json:"dgraph.type,omitempty"`
		EventID    string   `json:"event.id"`
		Summary    string   `json:"event.summary"`
		Timestamp  string   `json:"event.timestamp"`
		TrustScore float64  `json:"event.trust_score"`
	}

	// 1. Căutăm nodul existent
	query := `query q($id: string) {
		ev(func: eq(event.id, $id)) {
			uid
		}
	}`

	transaction := repo.graphClient.NewTxn()
	defer transaction.Discard(ctx)

	resp, err := transaction.QueryWithVars(ctx, query, map[string]string{"$id": eventID})
	if err != nil {
		return fmt.Errorf("failed to query event: %w", err)
	}

	var root struct {
		Ev []struct {
			Uid string `json:"uid"`
		} `json:"ev"`
	}
	if err := json.Unmarshal(resp.Json, &root); err != nil {
		return err
	}

	uid := "_:new"
	if len(root.Ev) > 0 {
		uid = root.Ev[0].Uid
	}

	// 2. Facem Upsert
	dto := EventDTO{
		Uid:        uid,
		DType:      []string{"Event"},
		EventID:    eventID,
		Summary:    summary,
		Timestamp:  timestamp.Format(time.RFC3339),
		TrustScore: score,
	}

	jsonData, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	mutation := &api.Mutation{
		SetJson:   jsonData,
		CommitNow: true,
	}

	_, err = transaction.Mutate(ctx, mutation)
	if err != nil {
		return fmt.Errorf("failed to upsert causal event: %w", err)
	}

	return nil
}
