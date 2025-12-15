package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/dgraph-io/dgo/v240"
	"github.com/dgraph-io/dgo/v240/protos/api"
	_ "github.com/lib/pq"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"google.golang.org/grpc"

	// "github.com/yourorg/truthweave/internal/infrastructure/arweave"
	"github.com/yourorg/truthweave/internal/infrastructure/dgraph"
	"github.com/yourorg/truthweave/internal/infrastructure/gemini"
	"github.com/yourorg/truthweave/internal/infrastructure/newsapi"
	"github.com/yourorg/truthweave/internal/infrastructure/postgres"

	// "github.com/yourorg/truthweave/internal/infrastructure/solana"
	"github.com/yourorg/truthweave/internal/infrastructure/temporal"
	"github.com/yourorg/truthweave/pkg/config"
)

// [RO] Punct de Intrare Lucrător (Worker)
//
// Acest program este "Muncitorul din Spate". El rulează pe un server separat (sau în container)
// și execută muncile grele (procesarea AI, scrierea în blockchain), fără să încetinească site-ul principal.
func main() {
	// [RO] 1. Configurare
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Eroare Worker: Nu am putut încărca configurările: %v", err)
	}

	// [RO] 2. Dependențe Infrastructură (Uneltele Muncitorului)

	// Postgres
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("Eroare Postgres: %v", err)
	}
	defer db.Close()
	pgRepo := postgres.NewPostgresNewsArticleRepository(db)

	// Dgraph
	dconn, err := grpc.Dial(cfg.DgraphHost, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Eroare Dgraph: %v", err)
	}
	dgraphClient := dgo.NewDgraphClient(api.NewDgraphClient(dconn))
	dgraphRepo := dgraph.NewDgraphKnowledgeGraphRepository(dgraphClient)

	// Google AI
	aiClient, err := gemini.NewGoogleGeminiArtificialIntelligenceAdapter(context.Background(), cfg.GeminiAPIKey)
	if err != nil {
		log.Fatalf("Eroare AI: %v", err)
	}

	// NewsAPI
	newsClient := newsapi.NewNewsAPIAdapter(cfg.NewsAPIKey)

	// Blockchain (Arweave & Solana) - DEZACTIVAT TEMPORAR PENTRU DEMO
	// Aici inițializăm adaptoarele reale. Dacă lipsesc cheile, worker-ul s-ar putea să crape la runtime.
	// Pentru demo, folosim inițializare "best effort" sau panicăm dacă e critic.
	// arweaveAdapter, _ := arweave.NewArweavePermanentStorageAdapter(cfg.ArweaveKeyPath, cfg.ArweaveGateway)

	// Solana necesită cheie privată. Dacă nu e în config, punem un string gol și adaptorul va returna eroare la constructor?
	// Constructorul returnează eroare. Vom gestiona eroarea.
	/*
		solanaAdapter, err := solana.NewSolanaBlockchainVerificationAdapter(cfg.SolanaPrivateKey, cfg.SolanaRPC)
		if err != nil {
			log.Printf("Avertisment: Modulul Solana nu a putut fi inițializat (cheie lipsă?): %v", err)
			// Putem continua, dar activitatea de ancorare va eșua (ceea ce e acceptabil pentru testare locală fară chei reale).
		}
	*/

	// [RO] 3. Conectare la Temporal Server
	tClient, err := client.Dial(client.Options{
		HostPort: cfg.TemporalHost,
	})
	if err != nil {
		log.Fatalf("Eroare Temporal: %v", err)
	}
	defer tClient.Close()

	// [RO] 4. Pornire Worker
	// "truthweave-task-queue" este canalul pe care ascultăm comenzi.
	w := worker.New(tClient, "truthweave-task-queue", worker.Options{})

	// Instanța care conține metodele ce vor fi executate
	activities := &temporal.NewsProcessingActivities{
		ArtificialIntelligence: aiClient,
		KnowledgeGraph:         dgraphRepo,
		Database:               pgRepo,
		NewsFetcher:            newsClient,
		DeduplicationThreshold: cfg.DeduplicationThreshold,
		// ArweaveArchiver:        arweaveAdapter,
		// BlockchainNotary:       solanaAdapter,
	}

	// Înregistrăm "Rețetele" (Flow-ul și Activitățile)
	w.RegisterWorkflow(temporal.OrchestrateNewsAnalysisWorkflow)
	w.RegisterWorkflow(temporal.GlobalNewsIngestionWorkflow)
	w.RegisterActivity(activities)

	log.Println("👷 Muncitorul TruthWeave este gata de treabă! Aștept comenzi...")
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("Muncitorul a întâmpinat o eroare fatală: %v", err)
	}
}
