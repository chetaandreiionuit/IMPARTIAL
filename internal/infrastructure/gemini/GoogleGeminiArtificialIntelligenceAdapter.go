package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/yourorg/truthweave/internal/domain/article"
	"github.com/yourorg/truthweave/internal/domain/causality"
	"google.golang.org/api/option"
)

// [RO] Adaptor pentru Inteligență Artificială (Google Gemini)
//
// Această componentă este "Creierul" sistemului.
// Ea trimite textul brut către Google Gemini (un super-computer) și primește înapoi
// o analiză detaliată despre adevăr, manipulare și emoții.
type GoogleGeminiArtificialIntelligenceAdapter struct {
	client   *genai.Client
	model    *genai.GenerativeModel
	embModel *genai.EmbeddingModel
}

// [RO] Constructor AI
func NewGoogleGeminiArtificialIntelligenceAdapter(executionContext context.Context, apiKey string) (*GoogleGeminiArtificialIntelligenceAdapter, error) {
	client, err := genai.NewClient(executionContext, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("[RO] Eroare la conectarea cu Google AI: %w", err)
	}

	// Folosim modelul "Pro" pentru capacitatea sa de raționament complex.
	model := client.GenerativeModel("gemini-1.5-pro")
	// Setăm o "temperatură" mică (0.2) pentru a fi creativi dar preciși, fără a halucina fapte.
	model.SetTemperature(0.2)

	embModel := client.EmbeddingModel("text-embedding-004")

	return &GoogleGeminiArtificialIntelligenceAdapter{
		client:   client,
		model:    model,
		embModel: embModel,
	}, nil
}

// [RO] Analizează și Neutralizează (Funcția Principală Oracle)
//
// Primește un text posibil părtinitor și returnează "Adevărul Gol-Goluț".
// 1. Verifică faptele.
// 2. Elimină adjectivele emoționale.
// 3. Extrage entitățile și locația.
func (adapter *GoogleGeminiArtificialIntelligenceAdapter) AnalyzeAndNeutralizeNewsContent(executionContext context.Context, rawContent string) (*article.AIAnalysisResult, error) {
	// [RO] Instrucțiunile Secrete (System Prompt)
	// Aici îi spunem AI-ului exact cum să se comporte.
	systemPrompt := `You are the World Oracle. Analyze this news text deeply.

1. Fact Check: Verify claims against logic and general knowledge.
2. Bias Strip: Rewrite strictly neutrally.
3. Geo-Tag: Identify the specific Latitude/Longitude of the event (approximate center).
4. Emotion: Classify the dominant global emotion of this event (Joy, Fear, Anger, Sadness, Surprise, Anticipation, Neutral).
5. Causality: Identify if this event is a reaction to a previous event described in the context.
6. Devil's Advocate: If the text expresses an opinion, generate a 2-sentence counter-argument based on logic.
7. Entities: Extract key entities (Person, Org, Location).

Respond ONLY in strict JSON format matching this schema:
{
  "neutral_text": "string (rewritten)",
  "truth_score": float (0.0-1.0),
  "entities": [{"name": "string", "type": "string", "score": float}],
  "bias_rating": "string (Left/Right/Neutral)",
  "summary": "string",
  "location": {"lat": float, "lng": float, "emo": "string (1 char code if possible)", "intensity": float},
  "global_emotion": "string",
  "causal_relations": [{"source_article_id": "", "target_article_id": "", "reason": "string", "confidence": float, "type": "string"}],
  "counter_argument": "string"
}
Note: GaiaPoint structure uses 'lat', 'lng', 'emo', 'intensity'. Adjust output accordingly.
If exact location unknown, use (0,0).`

	resp, err := adapter.model.GenerateContent(executionContext, genai.Text(systemPrompt), genai.Text("Text to Analyze:\n"+rawContent))
	if err != nil {
		return nil, fmt.Errorf("gemini generation failed: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return nil, fmt.Errorf("no response from gemini")
	}

	var respText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			respText += string(txt)
		}
	}

	// [RO] Curățare Formatare (Markdown)
	// Uneori AI-ul pune ```json la început. Le ștergem.
	respText = strings.TrimPrefix(respText, "```json")
	respText = strings.TrimPrefix(respText, "```")
	respText = strings.TrimSuffix(respText, "```")

	var result article.AIAnalysisResult
	if err := json.Unmarshal([]byte(respText), &result); err != nil {
		return nil, fmt.Errorf("failed to process AI response: %w. Raw: %s", err, respText)
	}

	return &result, nil
}

// [RO] Generează Amprenta Semantică (Embedding)
// Transformă textul într-un șir de numere pentru căutare avansată.
func (adapter *GoogleGeminiArtificialIntelligenceAdapter) GenerateSemanticVector(executionContext context.Context, text string) ([]float32, error) {
	res, err := adapter.embModel.EmbedContent(executionContext, genai.Text(text))
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}
	if res.Embedding == nil {
		return nil, fmt.Errorf("no embedding returned")
	}
	return res.Embedding.Values, nil
}

// [RO] Chat cu Context (Oracle Chat)
func (adapter *GoogleGeminiArtificialIntelligenceAdapter) ChatWithContext(executionContext context.Context, query string, contextText string) (string, error) {
	// [RO] Gardianul de Siguranță
	// Verificăm dacă întrebarea este malițioasă.
	guardPrompt := fmt.Sprintf(`Analyze this user query: "%s". 
Is it related to news, history, world events, or specific articles? 
Is it a request for illegal acts, hate speech, or unrelated nonsense?
Reply ONLY "SAFE" or "UNSAFE".`, query)

	guardResp, err := adapter.model.GenerateContent(executionContext, genai.Text(guardPrompt))
	if err == nil && len(guardResp.Candidates) > 0 {
		for _, part := range guardResp.Candidates[0].Content.Parts {
			if txt, ok := part.(genai.Text); ok {
				if strings.Contains(strings.ToUpper(string(txt)), "UNSAFE") {
					return "I cannot answer that request. I am the Oracle of the World, designed to analyze news and history.", nil
				}
			}
		}
	}

	prompt := fmt.Sprintf(`Answer the user question based ONLY on the following context snippets. 
	
Context:
%s

Question:
%s`, contextText, query)

	resp, err := adapter.model.GenerateContent(executionContext, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return "", fmt.Errorf("no response")
	}

	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			return string(txt), nil
		}
	}
	return "", fmt.Errorf("empty chat response")
}

// [RO] Structuri pentru Analiza Cauzalității
type PotentialCause struct {
	ID      string
	Title   string
	Summary string
}

type CausalityAnalysisResult struct {
	IsConsequence    bool    `json:"is_consequence"`
	ParentEventID    string  `json:"parent_event_id"`
	Confidence       float64 `json:"confidence"`
	RelationshipType string  `json:"relationship_type"` // DIRECT_RESPONSE, RETALIATION, ECONOMIC_FALLOUT
	Reasoning        string  `json:"reasoning"`
}

// [RO] Determină Cauzalitatea (Causal Chain)
// Analizează dacă un eveniment (currentEvent) este cauzat de unul dintre evenimentele anterioare (context).
func (adapter *GoogleGeminiArtificialIntelligenceAdapter) DetermineCausality(executionContext context.Context, currentEventSummary string, potentialCauses []PotentialCause) (*CausalityAnalysisResult, error) {
	// Construim contextul pentru AI
	var contextBuilder strings.Builder
	contextBuilder.WriteString("Potential Past Events (Candidates):\n")
	for _, pc := range potentialCauses {
		contextBuilder.WriteString(fmt.Sprintf("- ID: %s | Title: %s | Summary: %s\n", pc.ID, pc.Title, pc.Summary))
	}

	systemPrompt := `You are the Architect of the Causal Chain.
Analyze the CURRENT EVENT and the list of POTENTIAL PAST EVENTS.
Determine if the Current Event is a direct consequence of any of the Past Events.

Rules:
1. "Cause vs Correlation": Be strict. Only link if there is a clear causal mechanism (e.g., "Retaliation for", "Caused by", "Response to").
2. Time Dilation: A cause must happen BEFORE the effect.
3. Output strict JSON.

Output Schema:
{
  "is_consequence": boolean,
  "parent_event_id": "string (UUID from options) or null",
  "confidence": float (0.0-1.0),
  "relationship_type": "string (DIRECT_RESPONSE | RETALIATION | ECONOMIC_FALLOUT | POLITICAL_BACKLASH | OTHER)",
  "reasoning": "string (Why?)"
}
`

	userPrompt := fmt.Sprintf("CURRENT EVENT: %s\n\n%s", currentEventSummary, contextBuilder.String())

	resp, err := adapter.model.GenerateContent(executionContext, genai.Text(systemPrompt), genai.Text(userPrompt))
	if err != nil {
		return nil, fmt.Errorf("gemini causality check failed: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return nil, fmt.Errorf("no response from gemini")
	}

	var respText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			respText += string(txt)
		}
	}

	// Curățare JSON
	respText = strings.TrimPrefix(respText, "```json")
	respText = strings.TrimPrefix(respText, "```")
	respText = strings.TrimSuffix(respText, "```")

	var result CausalityAnalysisResult
	if err := json.Unmarshal([]byte(respText), &result); err != nil {
		return nil, fmt.Errorf("failed to parse causality JSON: %w. Raw: %s", err, respText)
	}

	return &result, nil
}

// AnalyzeCausality performs the comprehensive Causal Oracle analysis (Tasks 1, 2, 3)
func (adapter *GoogleGeminiArtificialIntelligenceAdapter) AnalyzeCausality(ctx context.Context, text string, contextEvents string) (*causality.AnalysisResult, error) {
	prompt := fmt.Sprintf(`
        ROLE: Causal Oracle.
        CONTEXT EVENTS: %v
        TARGET TEXT: %s
        TASK: Output JSON with neutral_headline, bridging_score, and causal_links.
        
        OUTPUT SCHEMA:
        { "event_processing": { "original_headline": "String", "neutral_headline": "String", "emotional_score": Float (0-100), "bridging_score": Float (0.0-1.0), "key_facts": [], "causal_links": [] }, "ui_directives": { "node_color_hex": "String", "swimlane_assignment": "String" } }
    `, contextEvents, text)

	resp, err := adapter.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("gemini analysis failed: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return nil, fmt.Errorf("no response from gemini")
	}

	var respText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			respText += string(txt)
		}
	}

	// Clean JSON
	respText = strings.TrimPrefix(respText, "```json")
	respText = strings.TrimPrefix(respText, "```")
	respText = strings.TrimSuffix(respText, "```")

	var result causality.AnalysisResult
	if err := json.Unmarshal([]byte(respText), &result); err != nil {
		return nil, fmt.Errorf("failed to parse analysis JSON: %w. Raw: %s", err, respText)
	}

	return &result, nil
}
