package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourorg/truthweave/internal/domain/ad"
)

// [RO] Manipulator Administrare Publicitate
//
// Gestionează panoul de control pentru reclame.
// Aici advertiserii pot crea și opri campanii.
type AdvertisementAdministrationHandlers struct {
	advertisementRepository ad.AdvertisementPersistenceInterface
}

// [RO] Constructor Admin
func NewAdvertisementAdministrationHandlers(repo ad.AdvertisementPersistenceInterface) *AdvertisementAdministrationHandlers {
	return &AdvertisementAdministrationHandlers{advertisementRepository: repo}
}

// [RO] Înregistrare Rute Admin
func (handler *AdvertisementAdministrationHandlers) RegisterAdminEndpoints(router *gin.Engine) {
	adminGroup := router.Group("/admin")
	{
		// [RO] POST /admin/ads -> Creează o reclamă nouă
		adminGroup.POST("/ads", handler.HandleCreateAdvertisementRequest)

		// [RO] PATCH /admin/ads/:id -> Activează/Dezactivează o reclamă
		adminGroup.PATCH("/ads/:id", handler.HandleToggleAdvertisementStatusRequest)
	}
}

// [RO] Manipulator: Creare Reclamă
func (handler *AdvertisementAdministrationHandlers) HandleCreateAdvertisementRequest(c *gin.Context) {
	var requestBody struct {
		Type      string `json:"type"`
		Title     string `json:"title"`
		Body      string `json:"body"`
		MediaURL  string `json:"media_url"`
		TargetURL string `json:"target_url"`
		Priority  int    `json:"priority"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON Invalid."})
		return
	}

	newCampaign := &ad.AdvertisementCampaignEntity{
		ID:        uuid.New(),
		Type:      requestBody.Type,
		Title:     requestBody.Title,
		Body:      requestBody.Body,
		MediaURL:  requestBody.MediaURL,
		TargetURL: requestBody.TargetURL,
		IsActive:  true,
		Priority:  requestBody.Priority,
		CreatedAt: time.Now(),
	}

	if err := handler.advertisementRepository.CreateAdvertisementCampaign(c.Request.Context(), newCampaign); err != nil {
		log.Printf("Eroare la crearea reclamei: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Nu am putut salva reclama."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": newCampaign.ID})
}

// [RO] Manipulator: Schimbare Status (On/Off)
func (handler *AdvertisementAdministrationHandlers) HandleToggleAdvertisementStatusRequest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Invalid."})
		return
	}

	var requestBody struct {
		Active bool `json:"active"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON Invalid."})
		return
	}

	if err := handler.advertisementRepository.SetAdvertisementActivationStatus(c.Request.Context(), id, requestBody.Active); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Nu am putut actualiza statusul."})
		return
	}
	c.Status(http.StatusOK)
}
