package ad

import (
	"time"

	"github.com/google/uuid"
)

// [RO] Entitate: Campanie Publicitară (Reclamă)
//
// Această structură definește o reclamă care poate fi afișată în aplicație.
// Modelul nostru de business este "Non-Intruziv": reclamele trebuie să fie relevante
// și clar marcate, fără să întrerupă lectura știrii.
type AdvertisementCampaignEntity struct {
	// [RO] Identificator Unic
	ID uuid.UUID `json:"id"`

	// [RO] Tipul Reclamei
	// "native" (se integrează în fluxul de știri) sau "banner" (clasic).
	Type string `json:"type"`

	// [RO] Titlu
	// Mesajul principal al reclamei (ex: "Descoperă noul telefon X").
	Title string `json:"title"`

	// [RO] Corpul Mesajului
	// Descrierea detaliată a produsului sau serviciului promovat.
	Body string `json:"body"`

	// [RO] Link Media (Imagine/Video)
	// URL către imaginea care va atrage atenția utilizatorului.
	MediaURL string `json:"media_url"`

	// [RO] Link Țintă
	// Unde ajunge utilizatorul dacă dă click (Website-ul advertiser-ului).
	TargetURL string `json:"target_url"`

	// [RO] Status Activ
	// Dacă este 'false', reclama nu va fi afișată nimănui (campanie oprită).
	IsActive bool `json:"is_active"`

	// [RO] Prioritate (Bidding)
	// Un număr care determină ordinea de afișare.
	// O prioritate mai mare înseamnă că reclama apare mai sus în listă.
	Priority int `json:"priority"`

	// [RO] Număr de Afișări (Impressions)
	// De câte ori a fost văzută această reclamă de către utilizatori.
	// Folosit pentru facturare și statistici.
	Impressions int64 `json:"impressions"`

	// [RO] Data Creării
	// Când a fost introdusă reclama în sistem.
	CreatedAt time.Time `json:"created_at"`
}
