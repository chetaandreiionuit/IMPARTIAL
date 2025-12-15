# 🌐 TRUTHWEAVE: The Global Impartial Oracle
> **Codename:** IMPARTIAL / ANTIGRAVITY  
> **Version:** 2.6 (Autonomous Intelligence Era)

## 📖 Introducere
TruthWeave este o platformă avansată de agregare și analiză a știrilor, proiectată pentru a elimina polarizarea și dezinformarea (Fake News) folosind Inteligența Artificială Generativă și analiza datelor la scară planetară. Sistemul nu doar "citește" știrile, ci le **înțelege**, le **verifică** încrucișat și le prezintă într-o interfață mobilă fluidă, guvernată de legile fizicii "Antigravity".

---

# �️ PART I: BACKEND ARCHITECTURE (The Brain)

Sistemul Backend este construit în **Go (Golang)** urmând cu strictețe **Clean Architecture** și **Domain-Driven Design (DDD)**. Este un sistem distribuit, orientat pe evenimente, orchestrat de **Temporal.io**.

## 1. Stack Tehnologic
*   **Limbaj:** Go 1.21+
*   **Web Framework:** Gin Gonic (Middleware, Routing)
*   **Orchestration:** Temporal.io (Workflow Engine)
*   **Databases:** 
    *   **PostgreSQL 15:** Persistență relațională, metadate, și (planned) pgvector.
    *   **Dgraph:** Knowledge Graph (Graful Cunoașterii) pentru conexiuni semantice între entități.
*   **AI:** Google Gemini 1.5 Pro (Generative Analysis, Neutralization).
*   **Scraping:** Colly V2 + Go-Readability (Autonomous Intel).
*   **Data Sources:** GDELT Project V2 DOC API, NewsAPI.
*   **Observability:** `log/slog` (Structured JSON Logging).

## 2. Structura Proiectului (Clean Architecture)
```text
/cmd
  /api          -> ApplicationEntryPoint.go (REST API Server)
  /worker       -> BackgroundWorkerEntryPoint.go (Temporal Worker)
/internal
  /domain       -> Reguli de Business PURE (Entități, Interfețe Repository). Nicio dependență externă.
  /usecase      -> Logica Aplicației (Services, Orchestrators). Leagă Domain de Infra.
  /infrastructure -> Implementări concrete (Postgres, Dgraph, Gemini, GDELT).
  /api          -> Stratul de prezentare HTTP (Handlers, Middleware, DTOs).
/pkg
  /config       -> Încărcare variabile de mediu (.env) via Viper.
  /logger       -> Sistem centralizat de logging structurat.
```

## 3. Fluxuri de Date și Clase Cheie

### A. Pipeline-ul Autonom de Ingestie (The "Sense" Phase)
Fluxul începe fără intervenția utilizatorului, declanșat de cron-joburi distribuite.

1.  **Trigger:** `GlobalNewsIngestionWorkflow` (Temporal Cron `*/15 * * * *`).
2.  **Discovery (GDELT):**
    *   Clasa: `internal/infrastructure/gdelt/GDELTAdapter.go`
    *   Funcționalitate: Construiește query-uri complexe (`toneabs>5`, `imagetag:tank`, `theme:CRISIS`) către API-ul GDELT V2.
    *   Output: O listă de `GdeltArticle` (metadate brute).
3.  **Extraction (Colly):**
    *   Clasa: `internal/infrastructure/scraper/CollyScraper.go`
    *   Tehnică: Utilizează `colly.Async(true)` și `LimitRule` (2 request-uri paralele/domeniu) pentru a evita blocarea.
    *   Curățare: HTML-ul brut este trecut prin `go-readability` direct din stream-ul HTTP (`OnResponse`), extrăgând doar textul semantic ("Clean content").

### B. Pipeline-ul de Analiză Cognitivă (The "Think" Phase)
Odată ce avem textul, intră în scenă "Creierul".

1.  **Orchestrator:** `NewsAnalysisWorkflowOrchestrator` (Temporal).
2.  **AI Processing (Gemini):**
    *   Clasa: `GoogleGeminiArtificialIntelligenceAdapter`
    *   Acțiune: Trimite textul curat către Gemini cu un "System Prompt" strict pentru neutralizare și fact-checking.
    *   Rezultat: `AIAnalysisResult` (Scor Adevăr 0-1, Bias Rating, Rezumat Imparțial).
3.  **Graph Construction (Dgraph):**
    *   Datele structurate sunt transformate în noduri (`Article`, `Person`, `Location`, `Organization`) și muchii (`MENTIONS`, `LOCATED_IN`) în graful Dgraph.

### C. Pipeline-ul de Livrare (The "Speak" Phase)
API-ul REST servește datele procesate.

1.  **Entry Point:** `ApplicationEntryPoint.go` -> inițializează `NewsArticleRequestHandlers`.
2.  **Middleware:** `internal/api/http/middleware/logging.go` -> Interceptează request-ul, generează `TraceID` și loghează în format JSON.
3.  **Handler:** `HandleFeedRequest` -> Apelează `NewsArticleOrchestrationService`.
    *   *Nota Bene:* Deși numele serviciului sugerează orchestrare, la citire el acționează ca un Façade peste Repository.

---

# 📱 PART II: ANDROID ARCHITECTURE (The Experience)

Aplicația Android este o demonstrație de forță tehnologică, implementând **Glassmorphism** și fizică avansată pe o arhitectură **MVI (Model-View-Intent)** rigidă.

## 1. Stack Tehnologic
*   **Limbaj:** Kotlin (JVM 1.8 compatibility).
*   **UI Framework:** Jetpack Compose (Material3).
*   **Architecture:** Clean Architecture + MVI.
*   **List Management:** Paging 3 (Infinite Scrolling).
*   **Dependency Injection:** Hilt (Dagger).
*   **Network:** Retrofit + OkHttp.
*   **Design System:** Custom "Antigravity" Theme.

## 2. Arhitectura MVI în Detaliu

### A. Fluxul Unidirecțional (UDF)
Datele curg într-o singură direcție, asigurând predictibilitatea stării.

`UI (Compose) -> Intent (Event) -> ViewModel -> UseCase -> Repository -> ViewModel (State/Flow) -> UI`

1.  **Intent:** Utilizatorul trage de listă ("Swipe to Refresh"). UI-ul emite `FeedIntent.Refresh`.
2.  **ViewModel** (`NewsFeedViewModel`):
    *   Interceptorul `handleIntent` primește evenimentul.
    *   Nu modifică starea direct. Apelează Repository-ul.
3.  **State Management:**
    *   Starea este împărțită în două:
        *   `pagingDataFlow`: `Flow<PagingData<FeedItem>>` pentru lista infinită. Este cache-uită în `viewModelScope` (`.cachedIn`) pentru a supraviețui rotației ecranului.
        *   `uiState`: Pentru erori globale sau loading inițial.

### B. Componente Cheie

1.  **NewsFeedScreen.kt (The View):**
    *   Observă `pagingDataFlow` folosind `collectAsLazyPagingItems()`.
    *   Folosește `LazyColumn` cu `contentType` pentru a randa eficient polimorfismul (Articole vs Reclame).
2.  **GlassModifier.kt (The Visuals):**
    *   Un `Modifier` custom care detectează versiunea de Android.
    *   **API 31+:** Aplică `RenderEffect.createBlurEffect` hardware-accelerated.
    *   **Legacy:** Aplică un fallback de transparență + noise texture.
3.  **FeedPagingSource.kt (The Data Pump):**
    *   Implementează logica de încărcare paginată.
    *   Funcția `getRefreshKey` calculează ancora matematică pentru a nu pierde poziția la invalidare.
4.  **NewsRepositoryImpl.kt (The Bridge):**
    *   Transformă apelurile API în `Pager` objects.

## 3. Design-ul "Antigravity"
*   **Physical Animation:** Elementele nu apar pur și simplu. Ele "intră" în scenă cu o animație `spring` (elastică) combinată cu `scaleIn` și `slideInVertically`, simulând imponderabilitatea.
*   **Spatial Gradients:** Fundalurile nu sunt culori solide, ci gradiente verticale complexe (`Deep Space Blue` -> `Nebula Violet`) care dau adâncime.

---

# 🛠️ PART III: SETUP & OPERATIONS

## 1. Pornirea Sistemului (Automatizată)
Am creat scriptul `start_backend.ps1` care:
1.  Lansează `docker-compose` (Postgres, Dgraph, Temporal).
2.  Așteaptă healthcheck-ul.
3.  Lansează `go run ./cmd/worker` (Procesare fundal).
4.  Lansează `go run ./cmd/api` (Server API).

Comandă:
```powershell
./start_backend.ps1
```

## 2. Compilarea Android
```powershell
cd android
./gradlew installDebug
```

## 3. Variabile de Mediu (.env)
Sistemul necesită un fișier `.env` în rădăcină. Configurația este încărcată via `pkg/config`.

```ini
SERVER_PORT=8080
DB_URL=postgres://user:password@localhost:5432/truthweave
GEMINI_API_KEY=AIza... (Cheia ta reală)
TEMPORAL_HOST=localhost:7233
```

---

# 📡 PART IV: API REFERENCE (Scurt extras)

Toate răspunsurile sunt JSON standardizat.

*   `GET /api/v1/news/feed?page=1&limit=10`
    *   Returnează fluxul principal (Articole + Reclame injectate).
*   `POST /api/v1/ingest`
    *   Trigger manual pentru analiză URL (`{"url": "..."}`).
*   `GET /api/v1/oracle/gaia-map`
    *   Date geospațiale pentru hartă.
*   `POST /api/v1/chat`
    *   Discuție cu agentul AI pe marginea unui articol.

---

> **Note finale:** Acest proiect reprezintă stadiul artei în ingineria software "Agentic", fiind scris și validat în proporție de 99% de agenți AI autonomi sub supraveghere umană.
