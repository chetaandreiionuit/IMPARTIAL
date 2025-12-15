# TRUTHWEAVE V2: THE CAUSAL ORACLE
### "De la Zgomot la Semnal. De la Haos la Cauzalitate."

---

## 📘 PARTEA 1: MANUALUL UTILIZATORULUI (NON-TEHNIC)

### Ce este TruthWeave?
TruthWeave nu este o aplicație de știri. Este un "Oracol Cauzal". Într-o lume inundată de clickbait și propagandă, TruthWeave funcționează ca un filtru de purificare a informației. Nu îți spune doar "ce s-a întâmplat", ci **de ce** s-a întâmplat și **ce urmează**.

### Cum funcționează? (Magia din Spate)
Imaginați-vă un analist veteran care citește tot internetul în timp real, elimină adjectivele inutile, verifică faptele din 3 surse contradictorii și desenează o hartă a evenimentelor pe o tablă. Asta face TruthWeave automat, folosind o rețea de Inteligență Artificială.

### Funcționalități Cheie (Features)

#### 1. The Causal Metro Map (Metroul Cauzalității)
Nu mai există liste infinite de articole. Știrile sunt afișate ca o **hartă de metrou verticală**.
*   **Linia Principală:** Firul narativ central al zilei (ex: "Criza Economică").
*   **Ramificații:** Când un eveniment declanșează altul (ex: "Inflație" -> declanșează -> "Proteste"), linia se bifurcă. Vedeți fizic legătura dintre cauză și efect.

#### 2. Emotional Noise Filter (Filtrul de Zgomot Emoțional)
Sunteți obosit de titluri care strigă la dumneavoastră? ("ȘOCANT! DEZASTRU!").
*   TruthWeave rescrie automat titlurile în timp real.
*   **Înainte:** *"Nu o să-ți vină să crezi! Piața crypto se prăbușește într-o baie de sânge!"*
*   **După (TruthWeave):** *"Bitcoin scade cu 5.2% pe fondul volumelor reduse de tranzacționare."*
*   Rezultatul este o experiență de lectură calmă, "Fintech", focusată pe date, nu pe dopamină.

#### 3. Bridging Consensus Score (Scorul de Adevăr Trans-Partinic)
Nu ne bazăm pe un singur "Fact Checker". Sistemul compară cum relatează aceeași știre surse opuse (ex: CNN vs Fox News vs Al Jazeera).
*   **100% Verified:** Știrea apare identic în toate sursele.
*   **Low Confidence:** Știrea apare doar într-o "bulă" informațională.

### Ghid de Utilizare Rapidă
1.  **Deschideți Aplicația:** Veți fi întâmpinat de "Metroul Cauzalității".
2.  **Navigare:** Scroll vertical pentru a merge în timp. Urmăriți liniile colorate pentru a vedea evoluția unui subiect.
3.  **Explorare:** Apăsați pe un nod (eveniment) pentru a vedea detaliile tehnice și scorul de încredere.

---

## 🛠 PARTEA 2: DOCUMENTAȚIA TEHNICĂ (INTERNAL BLUEPRINT)

### Arhitectura Sistemului
TruthWeave este un sistem distribuit, construit pe principiul **Event-Driven Architecture**.

**Core Stack:**
*   **Limbaj Backend:** Go 1.24 (Generics, High Performance)
*   **Orchestrare:** Temporal.io (Workflow Engine)
*   **Graph Database:** Dgraph (Stocarea relațiilor cauzale)
*   **AI Engine:** Google Gemini 1.5 Pro (Inference)
*   **Mobile App:** Android Native (Kotlin + Jetpack Compose)

### A. Fluxul Datelor (The Pipeline)

#### 1. Ingestie (GDELT Firehose)
*   **Clasa:** `GDELTAdapter.go` (`FetchHighImpactEvents`)
*   **Ce face:** Se conectează la stream-ul global GDELT V2 la fiecare 15 minute. Descarcă un CSV comprimat (~10MB), îl parsează în memorie (fără disc I/O greu) și filtrează evenimentele care au un "Tone Score" absolut mai mare de 5 (impact major).

#### 2. Orchestrare (Temporal)
*   **Clasa:** `CausalChainWorkflow.go`
*   **Ce face:** Este "dirijorul" infinit.
    *   Rulează într-o buclă `for` eternă.
    *   Primește un semnal `NewArticleSignal`.
    *   Execută secvențial activitățile: `Scrape` -> `AI Analysis` -> `Graph Upsert`.
    *   Gestionează erorile și retry-urile automat (dacă Gemini cade, Temporal reîncearcă exponential).

#### 3. Creierul (Gemini Analytics)
*   **Clasa:** `GoogleGeminiArtificialIntelligenceAdapter.go` (`AnalyzeCausality`)
*   **Ce face:**
    *   Primește textul brut al articolului.
    *   Injectează un "System Prompt" complex (The Causal Oracle).
    *   Cere un răspuns strict JSON care conține:
        *   `neutral_headline`: Titlul rescris.
        *   `bridging_score`: Scorul calculat.
        *   `causal_links`: O listă de ID-uri ale evenimentelor trecute care au cauzat acest eveniment.

#### 4. Memoria (Dgraph)
*   **Clasa:** `DgraphKnowledgeGraphRepository.go` (`UpsertCausalEvent`)
*   **Ce face:**
    *   Salvează entitatea `Event` în graf.
    *   Creează muchiile fizice (`<event_A> --[caused_by]--> <event_B>`).
    *   Aceasta este structura care permite vizualizarea "Metro Map".

---

## 📱 PARTEA 3: CLIENTUL ANDROID (FRONTEND)

### Arhitectura UI: "Editorial Fintech"
Aplicația este desenată procedural, fără imagini statice (png/jpg), pentru performanță și claritate maximă.

#### 1. `NewsFeedScreen.kt`
*   Containerul principal.
*   Folosește un fundal solid "Gunmetal" (`#121212`).

#### 2. `CausalGraphFeed.kt` & `CausalMetroMap.kt`
*   **Inima vizuală.**
*   Nu este o listă simplă (`RecyclerView`). Este un Canvas infinit.
*   Logica desenează linii Bezier între elementele din listă, bazat pe `parentID` din modelul de date.
*   Dacă un articol este "copilul" altuia, este indentat și conectat vizual.

#### 3. `EventNodeCard.kt`
*   Cardul individual de știre.
*   Design "Glassmorphism" pe fundal întunecat.
*   Afișează doar titlul neutru și scorul de încredere (verde/galben/roșu).

---

## 🚀 MANUAL DE PORNIRE (DEPLOYMENT)

### Cum pornesc proiectul local?

**Pasul 1: Infrastructura (Docker)**
Ai nevoie de Docker Desktop instalat.
Rulează în terminalul din rădăcina proiectului:
```powershell
docker-compose up -d
```
*Asta pornește Dgraph, Postgres și Temporal Server.*

**Pasul 2: Backend-ul (The Worker)**
Acesta este motorul care procesează datele.
```powershell
go run cmd/worker/BackgroundWorkerEntryPoint.go
```
*Vei vedea loguri cum că "Muncitorul este gata".*

**Pasul 3: Aplicația Mobilă**
1.  Deschide folderul `android` în **Android Studio**.
2.  Conectează un telefon sau pornește emulatorul.
3.  Apasă butonul **Play (Run)**.

### Cum testez?
Odată pornit worker-ul, el va începe automat să tragă date din GDELT. Dacă vrei să vezi "magia" instant:
1.  Uită-te în log-urile terminalului Worker. Vei vedea `[GDELT] Downloading update...` -> `Analysing with Gemini...` -> `Persisting Event`.
2.  Deschide aplicația Android. Vei vedea noile noduri apărând pe "hartă" în timp real (după refresh).
