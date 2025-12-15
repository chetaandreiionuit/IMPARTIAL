package ad

import (
	"context"

	"github.com/google/uuid"
)

// [RO] Interfața de Management a Reclamelor
//
// Definește operațiunile disponibile pentru gestionarea inventarului publicitar.
type AdvertisementPersistenceInterface interface {
	// [RO] Creează Campanie Nouă
	// Adaugă o nouă reclamă în sistem.
	CreateAdvertisementCampaign(execution_context context.Context, ad *AdvertisementCampaignEntity) error

	// [RO] Obține Reclamele Active
	// Returnează lista reclamelor care pot fi afișate acum, ordonate după prioritate.
	// De regulă, selectăm primele 5 cele mai valoroase reclame.
	RetrieveActiveAdvertisementCampaigns(execution_context context.Context) ([]*AdvertisementCampaignEntity, error)

	// [RO] Activează/Dezactivează Reclama
	// Pornește sau oprește afișarea unei campanii specifice.
	// (Ex: Advertiser-ul nu mai are buget -> setăm active = false).
	SetAdvertisementActivationStatus(execution_context context.Context, id uuid.UUID, isActive bool) error
}
