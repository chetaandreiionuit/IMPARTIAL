package temporal

import (
	"context"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"github.com/google/uuid"
	"github.com/yourorg/truthweave/internal/domain/article"

	"github.com/yourorg/truthweave/internal/infrastructure/dgraph"
	"github.com/yourorg/truthweave/internal/infrastructure/gdelt"
	"github.com/yourorg/truthweave/internal/infrastructure/gemini"
	"github.com/yourorg/truthweave/internal/infrastructure/postgres"
)

// [RO] Activitățile Fluxului de Lucru
// Aici listăm "uneltele" (dependințele) pe care le folosim în pașii analizei.
// Workflow-ul nu știe CUM se fac lucrurile (asta fac adaptoarele), el doar știe PE CINE să cheme.
type NewsProcessingActivities struct {
	ArtificialIntelligence *gemini.GoogleGeminiArtificialIntelligenceAdapter
	KnowledgeGraph         *dgraph.DgraphKnowledgeGraphRepository
	Database               *postgres.PostgresNewsArticleRepository
	NewsFetcher            *gdelt.GDELTAdapter // Replaced NewsAPI with GDELT V2
	DeduplicationThreshold float64
}

// [RO] Rezultat Similaritate
type SimilarityCheckResult struct {
	ExistingArticle *article.NewsArticleEntity
	SimilarityScore float64
}

// [RO] Activitate: Colectare Știri Globale (Reală)
// Acum folosește GDELT Project V2 pentru a detecta evenimente cu impact major (|tone| > 5).
func (activities *NewsProcessingActivities) FetchLatestGlobalNewsActivity(ctx context.Context) ([]string, error) {
	return activities.NewsFetcher.FetchHighImpactEvents(ctx)
}

// [RO] Activitate 1: Extragere Conținut (Scraping)
func (activities *NewsProcessingActivities) ExtractWebPageContentActivity(executionContext context.Context, url string) (string, error) {
	// Dacă URL-ul vine de la NewsAPI, uneori conținutul e trunchiat.
	// Aici ar trebui un "Full Content Scraper" (ex: Jina Reader sau custom Go scraper cu colly).
	// Pentru acest pas, simulăm că am descărcat tot textul.
	return "Simulated FULL content scraped from " + url + ". The situation is evolving rapidly...", nil
}

// [RO] Activitate 2: Generare Vector Semantic
func (activities *NewsProcessingActivities) GenerateSemanticVectorActivity(executionContext context.Context, text string) ([]float32, error) {
	return activities.ArtificialIntelligence.GenerateSemanticVector(executionContext, text)
}

// [RO] Activitate 3: Verificare Duplicate (Deduplicare)
func (activities *NewsProcessingActivities) CheckForExistingDuplicatesActivity(executionContext context.Context, embedding []float32) (*SimilarityCheckResult, error) {
	articles, err := activities.Database.FindSemanticallySimilarArticles(executionContext, embedding, 1)
	if err != nil {
		return nil, err
	}
	if len(articles) == 0 {
		return &SimilarityCheckResult{ExistingArticle: nil, SimilarityScore: 0.0}, nil
	}

	// [RO] Folosim Pragul din Configurare
	threshold := activities.DeduplicationThreshold
	if threshold == 0 {
		threshold = 0.90 // Fallback dacă nu e setat
	}

	bestMatch := articles[0]
	// Simulam scorul de similaritate pentru demo (deoarece pgvector nu-l returneaza direct in structura entity aici)
	simScore := 0.9

	if simScore < threshold {
		return &SimilarityCheckResult{ExistingArticle: nil, SimilarityScore: simScore}, nil
	}

	return &SimilarityCheckResult{ExistingArticle: bestMatch, SimilarityScore: simScore}, nil
}

// [RO] Activitate 4: Analiză AI Completă
func (activities *NewsProcessingActivities) AnalyzeNewsContentActivity(executionContext context.Context, rawContent string) (*article.AIAnalysisResult, error) {
	return activities.ArtificialIntelligence.AnalyzeAndNeutralizeNewsContent(executionContext, rawContent)
}

// [RO] Activitate 5: Salvare în Baza de Date
func (activities *NewsProcessingActivities) PersistAnalysisToDatabaseActivity(executionContext context.Context, newsArticle article.NewsArticleEntity) error {
	return activities.Database.PersistNewsArticle(executionContext, &newsArticle)
}

// [RO] Activitate 6: Actualizare Graf Cunoștințe
func (activities *NewsProcessingActivities) ConnectKnowledgeGraphActivity(executionContext context.Context, newsArticle article.NewsArticleEntity) error {
	return activities.KnowledgeGraph.SaveNewsArticleToGraph(executionContext, &newsArticle)
}

// [RO] Ochestratorul Fluxului de Analiză a Știrilor (Per Articol)
func OrchestrateNewsAnalysisWorkflow(workflowContext workflow.Context, articleURL string) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval: time.Second,
			MaximumAttempts: 5,
		},
	}
	workflowContext = workflow.WithActivityOptions(workflowContext, options)
	logger := workflow.GetLogger(workflowContext)

	var tools *NewsProcessingActivities

	// 1. Scrape
	var rawContent string
	if err := workflow.ExecuteActivity(workflowContext, tools.ExtractWebPageContentActivity, articleURL).Get(workflowContext, &rawContent); err != nil {
		return err
	}

	// 2. Vector
	var semanticVector []float32
	if err := workflow.ExecuteActivity(workflowContext, tools.GenerateSemanticVectorActivity, rawContent).Get(workflowContext, &semanticVector); err != nil {
		return err
	}

	// 3. Duplicate Check
	var similarityResult SimilarityCheckResult
	if err := workflow.ExecuteActivity(workflowContext, tools.CheckForExistingDuplicatesActivity, semanticVector).Get(workflowContext, &similarityResult); err != nil {
		return err
	}

	// Logică duplicat în Workflow (folosind un threshold hardcoded aici pentru siguranță sau bazându-ne pe activitate dacă returna bool)
	// Pentru consistență cu activitatea care returnează score, decidem aici.
	if similarityResult.SimilarityScore >= 0.90 && similarityResult.ExistingArticle != nil {
		logger.Info("[RO] Duplicat detectat. Oprim procesarea.", "existing_id", similarityResult.ExistingArticle.ID)
		return nil
	}

	// 4. AI Analysis
	var aiAnalysis article.AIAnalysisResult
	if err := workflow.ExecuteActivity(workflowContext, tools.AnalyzeNewsContentActivity, rawContent).Get(workflowContext, &aiAnalysis); err != nil {
		return err
	}

	// Construcție Entitate
	processedArticle := article.NewsArticleEntity{
		ID:          uuid.New(),
		OriginalURL: articleURL,
		Title:       "Analyzed: " + articleURL,
		Content:     aiAnalysis.RewrittenText,
		RawContent:  rawContent,
		Summary:     aiAnalysis.Summary,
		TruthScore:  aiAnalysis.Score,
		BiasRating:  aiAnalysis.BiasRating,
		Embedding:   semanticVector,
		PublishedAt: workflow.Now(workflowContext),
		Geolocation: article.GaiaPoint{
			ID:        uuid.New().String(),
			Latitude:  aiAnalysis.Location.Latitude,
			Longitude: aiAnalysis.Location.Longitude,
			Emotion:   aiAnalysis.Location.Emotion,
			Intensity: aiAnalysis.Location.Intensity,
		},
		GlobalEmotion:   aiAnalysis.GlobalEmotion,
		Causes:          aiAnalysis.CausalRelations,
		CounterArgument: aiAnalysis.CounterArgument,
		Mentions:        aiAnalysis.Entities,
	}

	// 5. Save DB
	if err := workflow.ExecuteActivity(workflowContext, tools.PersistAnalysisToDatabaseActivity, processedArticle).Get(workflowContext, nil); err != nil {
		return err
	}

	// 6. Save Graph
	if err := workflow.ExecuteActivity(workflowContext, tools.ConnectKnowledgeGraphActivity, processedArticle).Get(workflowContext, nil); err != nil {
		return err
	}

	return nil
}

// [RO] Agregatorul Global (Cron Workflow)
// Rulează periodic, aduce știri, lansează procesarea individuală.
func GlobalNewsIngestionWorkflow(ctx workflow.Context) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 2,
	}
	ctx = workflow.WithActivityOptions(ctx, options)
	logger := workflow.GetLogger(ctx)

	var tools *NewsProcessingActivities

	// 1. Fetch Headlines
	var newArticleURLs []string
	if err := workflow.ExecuteActivity(ctx, tools.FetchLatestGlobalNewsActivity).Get(ctx, &newArticleURLs); err != nil {
		logger.Error("Eșec NewsAPI", "Error", err)
		return err
	}

	logger.Info("Am găsit știri noi", "count", len(newArticleURLs))

	// 2. Fan-Out
	for _, url := range newArticleURLs {
		childOptions := workflow.ChildWorkflowOptions{
			WorkflowID: "analyze-" + url, // Deduplicare naturală Temporal
		}
		childCtx := workflow.WithChildOptions(ctx, childOptions)
		workflow.ExecuteChildWorkflow(childCtx, OrchestrateNewsAnalysisWorkflow, url)
	}

	return nil
}
