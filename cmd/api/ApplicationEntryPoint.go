package main

import (
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.temporal.io/sdk/client"

	server "github.com/yourorg/truthweave/internal/api/http"
	"github.com/yourorg/truthweave/internal/api/http/middleware"
	"github.com/yourorg/truthweave/internal/infrastructure/gemini"
	"github.com/yourorg/truthweave/internal/infrastructure/postgres"
	"github.com/yourorg/truthweave/internal/infrastructure/temporal"
	"github.com/yourorg/truthweave/internal/usecase/article"
	"github.com/yourorg/truthweave/pkg/config"
	"github.com/yourorg/truthweave/pkg/logger"
)

// [RO] Punctul de Intrare în Aplicație (Main)
//
// Aici începe totul. Această funcție este "Constructorul Suprem" care asamblează toate piesele LEGO:
// 1. Citește Configurația.
// 2. Conectează Bazele de Date (Postgres, Dgraph).
// 3. Inițializează Serviciile Externe (AI, Temporal).
// 4. Creează Managerii Logici (Services).
// 5. Pornește Serverul Web (API).
func main() {
	// [RO] 0. Inițializare Logger Centralizat
	appLogger := logger.InitLogger(logger.Config{
		ServiceName: "truthweave-api",
		Environment: "development",
		Level:       "info",
	})
	appLogger.Info("Sistemul TruthWeave se inițializează...")

	// [RO] 1. Încărcare Configurație
	cfg, err := config.LoadConfig()
	if err != nil {
		appLogger.Error("Eroare Critică: Nu am putut citi setările (config)", "error", err)
		return
	}

	// [RO] 2. Conectare la PostgreSQL (Memoria de Lungă Durată)
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		appLogger.Error("Eroare Critică: Conexiunea la Postgres a eșuat", "error", err)
		return
	}
	defer db.Close()

	// [RO] 3. Inițializare Depozite (Repositories)
	// Creăm "Bibliotecarii" care se ocupă de date.
	newsRepository := postgres.NewPostgresNewsArticleRepository(db)
	adRepository := postgres.NewPostgresAdvertisementRepository(db)

	// [RO] 4. Conectare la Temporal (Orchestratorul de Procese)
	tClient, err := client.Dial(client.Options{
		HostPort: cfg.TemporalHost,
	})
	if err != nil {
		appLogger.Error("Eroare Critică: Nu mă pot conecta la Temporal", "error", err)
		return
	}
	defer tClient.Close()
	temporalOrchestrator := temporal.NewTemporalOrchestratorClient(tClient)

	// [RO] 5. Conectare la Google Gemini (Creierul AI)
	aiClient, err := gemini.NewGoogleGeminiArtificialIntelligenceAdapter(context.Background(), cfg.GeminiAPIKey)
	if err != nil {
		appLogger.Error("Eroare Critică: AI-ul nu răspunde", "error", err)
		return
	}

	// [RO] 6. Asamblare Serviciu Principal (Business Logic)
	// Aici injectăm toate dependințele în "Managerul" aplicației.
	newsService := article.NewNewsArticleOrchestrationService(
		newsRepository,
		adRepository,
		temporalOrchestrator,
		aiClient,
	)

	// [RO] 7. Configurare Controller HTTP (API)
	// Pregătim "Recepția" care va răspunde la cererile mobile.
	httpHandler := server.NewNewsArticleRequestHandlers(newsService)
	adminHandler := server.NewAdvertisementAdministrationHandlers(adRepository)

	// [RO] 8. Start Server (Cu Middleware Logger)
	r := gin.New()
	r.Use(gin.Recovery()) // Panic recovery standard
	r.Use(middleware.StructuredLogger(appLogger))

	httpHandler.RegisterAPIEndpoints(r)
	adminHandler.RegisterAdminEndpoints(r)

	appLogger.Info("🚀 Aplicația TruthWeave a pornit cu succes!", "port", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		appLogger.Error("Serverul s-a oprit neașteptat", "error", err)
	}
}
