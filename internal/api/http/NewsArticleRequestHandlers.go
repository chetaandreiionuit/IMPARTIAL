package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/truthweave/internal/usecase/article"
)

// [RO] Manipulator Cereri HTTP (Controller)
//
// Această componentă este "Recepția" hotelului nostru.
// Ea întâmpină clienții (aplicația mobilă), le ascultă cererile, le validează
// și le trimite mai departe către Manager (Service).
// La final, tot ea le dă răspunsul (JSON).
type NewsArticleRequestHandlers struct {
	orchestrationService *article.NewsArticleOrchestrationService
}

// [RO] Constructor Controller
func NewNewsArticleRequestHandlers(service *article.NewsArticleOrchestrationService) *NewsArticleRequestHandlers {
	return &NewsArticleRequestHandlers{orchestrationService: service}
}

// [RO] Înregistrare Rute (Harta API)
// Definește URL-urile la care aplicația răspunde.
func (handler *NewsArticleRequestHandlers) RegisterAPIEndpoints(router *gin.Engine) {
	apiGroup := router.Group("/api/v1")
	{
		// [RO] POST /ingest -> Trimite un link pentru analiză
		apiGroup.POST("/ingest", handler.HandleIngestionRequest)

		// [RO] GET /news/:id -> Citește o știre analizată
		apiGroup.GET("/news/:id", handler.HandleGetNewsRequest)

		// [RO] GET /news/feed -> Obține fluxul de noutăți (cu reclame)
		apiGroup.GET("/news/feed", handler.HandleFeedRequest)

		// [RO] POST /chat -> Vorbește cu Oracolul
		apiGroup.POST("/chat", handler.HandleOracleChatRequest)

		// [RO] GET /oracle/gaia-map -> Harta Adevărului
		apiGroup.GET("/oracle/gaia-map", handler.HandleGaiaMapRequest)
	}
}

// [RO] Manipulator: Ingestie (Procesare)
func (handler *NewsArticleRequestHandlers) HandleIngestionRequest(c *gin.Context) {
	var requestBody struct {
		URL string `json:"url"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cerere Invalidă. Avem nevoie de un URL."})
		return
	}

	response, err := handler.orchestrationService.StartNewsAnalysisPipeline(c.Request.Context(), article.NewsProcessingRequest{TargetURL: requestBody.URL})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// [RO] Răspuns: 202 Accepted
	// Spunem "Am primit comanda, lucrăm la ea".
	c.JSON(http.StatusAccepted, gin.H{"job_id": response.JobID})
}

// [RO] Manipulator: Citire Știre
func (handler *NewsArticleRequestHandlers) HandleGetNewsRequest(c *gin.Context) {
	id := c.Param("id")
	newsArticle, err := handler.orchestrationService.RetrieveCompleteNewsArticle(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Articolul nu a fost găsit."})
		return
	}

	// [RO] Răspuns JSON pentru Mobil
	c.JSON(http.StatusOK, gin.H{
		"article": gin.H{
			"id":       newsArticle.ID,
			"headline": newsArticle.Title,
			"body":     newsArticle.Content,
			"truth_stats": gin.H{
				"score":          newsArticle.TruthScore,
				"bias_direction": newsArticle.BiasRating,
			},
			"updated_at": newsArticle.ProcessedAt,
		},
	})
}

// [RO] Manipulator: Flux de Știri
func (handler *NewsArticleRequestHandlers) HandleFeedRequest(c *gin.Context) {
	feed, err := handler.orchestrationService.GeneratePersonalizedNewsFeed(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"articles": feed})
}

// [RO] Manipulator: Chat Oracle
func (handler *NewsArticleRequestHandlers) HandleOracleChatRequest(c *gin.Context) {
	var requestBody struct {
		ArticleID string `json:"article_id"`
		Question  string `json:"question"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format invalid."})
		return
	}

	answer, citations, err := handler.orchestrationService.PerformOracleContextualSearch(c.Request.Context(), requestBody.Question, requestBody.ArticleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"answer":    answer,
		"citations": citations,
	})
}

// [RO] Manipulator: Harta Gaia
func (handler *NewsArticleRequestHandlers) HandleGaiaMapRequest(c *gin.Context) {
	points, err := handler.orchestrationService.AggregateGlobalTruthMap(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, points)
}
