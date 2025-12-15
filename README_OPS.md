# ⚙️ Ghid Operațional: TruthWeave (TruthWeave Ops)

Acest document descrie procedurile operaționale standard pentru rularea, configurarea și administrarea platformei **TruthWeave**.

---

## 🔧 Configurare Avansată

Sistemul utilizează variabile de mediu pentru ajustarea parametrilor cheie, fără a necesita recompilare.

### 1. Pragul de Deduplicare AI (`DEDUPLICATION_THRESHOLD`)

Controlează cât de "strică" este Inteligența Artificială când decide dacă două știri vorbesc despre același lucru.

*   **Valoare Implicită:** `0.90` (90% similaritate)
*   **Cum funcționează:**
    *   O valoare **mai mică** (ex: `0.80`) va grupa mai agresiv știrile (risc: poate grupa știri distincte dar similare).
    *   O valoare **mai mare** (ex: `0.98`) va necesita ca știrile să fie aproape identice pentru a fi considerate duplicate.

**Setare via Environment Variable:**
```bash
export DEDUPLICATION_THRESHOLD=0.92
```

---

## 🚀 Pornirea Sistemului (Docker)

Pentru a porni întreaga stivă (Bază de date, AI Worker, API, Dgraph) într-un mediu de producție sau staging:

```bash
docker-compose up -d --build
```

### Verificare Sănătate (Health Checks)

*   **API Server:** `http://localhost:8080/health` (trebuie să răspundă cu `200 OK`)
*   **Temporal UI:** `http://localhost:8088` (pentru monitorizarea fluxurilor AI)

---

## 👮 Administrare și Primul Utilizator

Deoarece sistemul este descentralizat, nu există un "Super Admin" hardcodat. Primul administrator trebuie injectat direct în baza de date.

### SQL pentru injectarea Admin-ului Suprem:

Rulați acest query în consola PostgreSQL:

```sql
INSERT INTO users (id, email, role, created_at)
VALUES (
    '00000000-0000-0000-0000-000000000000', 
    'admin@truthweave.org', 
    'ADMIN', 
    NOW()
);
```

---

## 🌍 Politica "No Null Island" (Geo-Sanitizare)

Pentru vizualizarea 3D (Gaia Map), sistemul aplică automat următoarele reguli:

1.  **Excludere:** Punctele care au coordonatele `(0.0, 0.0)` sau `NULL` sunt complet excluse din API-ul de hartă.
2.  **Fallback (În Dezvoltare):** Dacă AI-ul detectează o țară dar nu un oraș, va folosi centroidul țării respective.

---

## 🧹 Mentenanță Periodică

*   **Curățare Log-uri:** Docker log-urile trebuie rotite la fiecare 7 zile.
*   **Re-indexare Vectorială:** O dată pe lună, rulați `REINDEX INDEX articles_embedding_idx;` în Postgres pentru a menține performanța căutărilor semantice.
