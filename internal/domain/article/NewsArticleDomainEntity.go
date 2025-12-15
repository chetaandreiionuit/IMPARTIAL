package article

import (
	"errors"
	"net/url"
	"time"

	"github.com/google/uuid"
)

// [RO] Entitate de Domeniu: Știre (Articol de Presă)
//
// Această structură reprezintă nucleul aplicației noastre. Ea definește ce înseamnă
// o "Știre" în sistemul TruthWeave. Nu este doar un text, ci un obiect complex
// care conține analiza de adevăr, emoția globală și legăturile cauzale detectate.
//
// Gândiți-vă la această entitate ca la un "Dosar de Investigație" pentru fiecare articol
// citit de pe internet.
type NewsArticleEntity struct {
	// [RO] Identificator Unic (ID)
	// Un cod unic generat automat (de tip UUID) care ne ajută să găsim acest articol
	// în baza de date, indiferent dacă titlul sau url-ul se schimbă.
	ID uuid.UUID `json:"id"`

	// [RO] Adresa Web Originală (URL)
	// Link-ul exact de unde a fost preluat articolul. Este folosit pentru a verifica
	// dacă am mai procesat acest articol (deduplicare).
	OriginalURL string `json:"original_url"`

	// [RO] Titlul Știrii
	// Titlul extras din pagina web originală.
	Title string `json:"title"`

	// [RO] Conținutul Rescris (Neutralizat)
	// Aici stocăm varianta "curată" a știrii. Inteligența Artificială a rescris
	// textul original pentru a elimina opiniile subiective și limbajul inflamator,
	// păstrând doar faptele brute.
	Content string `json:"content"`

	// [RO] Conținutul Brut
	// Textul exact așa cum a fost găsit pe site-ul sursă, păstrat pentru arhivare
	// și comparație legală.
	RawContent string `json:"raw_content,omitempty"`

	// [RO] Rezumat Executiv
	// O sinteză scurtă (2-3 fraze) generată de AI pentru a fi citită rapid pe mobil.
	Summary string `json:"summary"`

	// [RO] Scorul de Adevăr (Trust Score)
	// O notă de la 0.0 la 1.0 dată de Algoritmul Oracle.
	// 1.0 înseamnă "Fapte Verificate 100%", 0.0 înseamnă "Fals Total / Propagandă".
	TruthScore float64 `json:"truth_score"`

	// [RO] Rating de Părtinire (Bias)
	// O etichetă care descrie orientarea politică sau emoțională (ex: "Left-Leaning", "Neutral", "Pro-Gov").
	BiasRating string `json:"bias_rating"`

	// [RO] Amprenta Semantică (Vector)
	// O listă de numere care reprezintă "înțelesul" matematic al textului.
	// Este folosită pentru a găsi alte articole similare în baza de date vectorială.
	Embedding []float32 `json:"-"`

	// [RO] Data Publicării
	// Când a apărut știrea pe internet.
	PublishedAt time.Time `json:"published_at"`

	// [RO] Data Procesării
	// Când a analizat sistemul nostru această știre.
	ProcessedAt time.Time `json:"processed_at"`

	// --- Secțiunea Verificare Blockchain ---

	// [RO] Identificator Arhivă Arweave
	// Dovada că acest articol a fost salvat permanent pe "hard disk-ul etern" al internetului (Arweave).
	// Nimeni nu îl mai poate șterge sau modifica.
	ArweaveTransactionID string `json:"arweave_tx_id,omitempty"`

	// [RO] Semnătura Digitală Solana
	// O "ștampilă" criptografică pe blockchain-ul Solana care atestă că noi (TruthWeave)
	// am validat acest articol la o anumită dată.
	SolanaSignature string `json:"solana_signature,omitempty"`

	// --- Secțiunea Oracle Genesis (Analiză Avansată) ---

	// [RO] Locația Geografică (Gaia Point)
	// Coordonatele GPS asociate principalului eveniment din știre, pentru afișarea pe Globul 3D.
	Geolocation GaiaPoint `json:"geo_location"`

	// [RO] Emoția Globală
	// Sentimentul general detectat în text (ex: "Frică", "Bucurie", "Furie").
	GlobalEmotion string `json:"global_emotion"`

	// [RO] Cauze (De ce s-a întâmplat?)
	// Lista de evenimente anterioare care au dus la această știre.
	Causes []CausalEventLink `json:"causes,omitempty"`

	// [RO] Efecte (Ce urmează?)
	// Lista de posibile consecințe prezise de AI.
	Effects []CausalEventLink `json:"effects,omitempty"`

	// [RO] Oglinda Adevărului (Counter-Argument)
	// Un argument automat generat care oferă perspectiva opusă, pentru a sparge "bula" cititorului.
	CounterArgument string `json:"counter_argument,omitempty"`

	// [RO] Entități Menționate
	// Lista de persoane, organizații sau locuri detectate în text.
	Mentions []NamedEntity `json:"mentions,omitempty"`
}

// [RO] Punct Geografic (Gaia)
// Reprezintă un punct pe harta 3D (Globul Adevărului).
type GaiaPoint struct {
	ID        string  `json:"id"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Emotion   string  `json:"emo"`       // Cod Emoție (ex: 'F' - Fear)
	Intensity float64 `json:"intensity"` // Intensitatea (0.0 - 1.0)
}

// [RO] Entitate Numită (Persoană/Org)
// Cineva sau ceva important menționat în știre.
type NamedEntity struct {
	Name  string  `json:"name"`
	Type  string  `json:"type"` // Person, Organization, Place
	Score float64 `json:"score"`
}

// [RO] Legătură Cauzală
// O conexiune logică între două articole (ex: "Articolul B este reacția la Articolul A").
type CausalEventLink struct {
	ID              string  `json:"id,omitempty"`
	SourceArticleID string  `json:"source_article_id"`
	TargetArticleID string  `json:"target_article_id"`
	Reason          string  `json:"reason"` // Explicația legăturii
	Confidence      float64 `json:"confidence"`
	Type            string  `json:"type"` // "caused_by" (cauzat de) sau "triggered" (declanșat)
}

// [RO] Rezultatul Analizei AI
// Structura care primește datele brute de la Gemini AI.
type AIAnalysisResult struct {
	RewrittenText   string            `json:"neutral_text"`
	Score           float64           `json:"truth_score"`
	Entities        []NamedEntity     `json:"entities"`
	BiasRating      string            `json:"bias_rating"`
	Summary         string            `json:"summary"`
	Location        GaiaPoint         `json:"location"`
	GlobalEmotion   string            `json:"global_emotion"`
	CausalRelations []CausalEventLink `json:"causal_relations"`
	CounterArgument string            `json:"counter_argument"`
}

// [RO] Sentiment AI
// O structură auxiliară pentru analiza de sentiment.
type AISentimentAnalysis struct {
	TruthScore      float64       `json:"truth_score"`
	BiasRating      string        `json:"bias_rating"`
	GlobalEmotion   string        `json:"global_emotion"`
	CounterArgument string        `json:"counter_argument"`
	Entities        []NamedEntity `json:"entities"`
}

// [RO] Validarea Datelor (Business Logic)
//
// Această metodă verifică dacă știrea respectă regulile minime ale sistemului
// înainte de a fi procesată. De exemplu, o știre fără URL este invalidă.
func (article *NewsArticleEntity) VerifyDataIntegrity() error {
	if article.OriginalURL == "" {
		return errors.New("[RO] Eroare Critică: Articolul nu are un URL sursă valid.")
	}
	if _, err := url.ParseRequestURI(article.OriginalURL); err != nil {
		return errors.New("[RO] Eroare Critică: Formatul URL-ului este incorect.")
	}
	return nil
}
