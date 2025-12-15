package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/yourorg/truthweave/internal/domain/ad"
)

// [RO] Depozit de Date PostgreSQL pentru Reclame
//
// Gestionează campaniile publicitare ("Ads") în baza de date.
// Asigură că reclamele sunt servite rapid și că bugetele sunt respectate.
type PostgresAdvertisementRepository struct {
	databaseConnection *sql.DB
}

// [RO] Constructor pentru Reclame
func NewPostgresAdvertisementRepository(db *sql.DB) *PostgresAdvertisementRepository {
	return &PostgresAdvertisementRepository{databaseConnection: db}
}

// [RO] Creează Campanie (Implementare)
func (repo *PostgresAdvertisementRepository) CreateAdvertisementCampaign(executionContext context.Context, advertisement *ad.AdvertisementCampaignEntity) error {
	sqlQuery := `
		INSERT INTO ads (id, type, title, body, media_url, target_url, is_active, priority, impressions, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 0, $9)
	`
	_, err := repo.databaseConnection.ExecContext(executionContext, sqlQuery,
		advertisement.ID,
		advertisement.Type,
		advertisement.Title,
		advertisement.Body,
		advertisement.MediaURL,
		advertisement.TargetURL,
		advertisement.IsActive,
		advertisement.Priority,
		time.Now(),
	)
	return err
}

// [RO] Obține Reclame Active (Implementare)
func (repo *PostgresAdvertisementRepository) RetrieveActiveAdvertisementCampaigns(executionContext context.Context) ([]*ad.AdvertisementCampaignEntity, error) {
	sqlQuery := `
		SELECT id, type, title, body, media_url, target_url,  priority, impressions
		FROM ads WHERE is_active = true ORDER BY priority DESC LIMIT 5
	`
	rows, err := repo.databaseConnection.QueryContext(executionContext, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activeAds []*ad.AdvertisementCampaignEntity
	for rows.Next() {
		var currentAd ad.AdvertisementCampaignEntity
		if err := rows.Scan(
			&currentAd.ID,
			&currentAd.Type,
			&currentAd.Title,
			&currentAd.Body,
			&currentAd.MediaURL,
			&currentAd.TargetURL,
			&currentAd.Priority,
			&currentAd.Impressions,
		); err != nil {
			continue
		}
		currentAd.IsActive = true
		activeAds = append(activeAds, &currentAd)
	}
	return activeAds, nil
}

// [RO] Activează/Dezactivează Reclama (Implementare)
func (repo *PostgresAdvertisementRepository) SetAdvertisementActivationStatus(executionContext context.Context, id uuid.UUID, isActive bool) error {
	sqlQuery := "UPDATE ads SET is_active = $1 WHERE id = $2"
	_, err := repo.databaseConnection.ExecContext(executionContext, sqlQuery, isActive, id)
	return err
}
