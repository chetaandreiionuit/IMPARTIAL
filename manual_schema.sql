-- 001_initial_schema.up.sql
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS articles (
    id UUID PRIMARY KEY,
    original_url TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    raw_content TEXT NOT NULL DEFAULT '',
    summary TEXT NOT NULL DEFAULT '',
    truth_score DOUBLE PRECISION NOT NULL,
    bias_rating TEXT NOT NULL,
    embedding vector(768), -- Gemini 1.5/Gecko embedding size
    published_at TIMESTAMPTZ NOT NULL,
    processed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS articles_embedding_idx ON articles USING hnsw (embedding vector_cosine_ops);

-- 002_oracle_schema.up.sql
ALTER TABLE articles ADD COLUMN IF NOT EXISTS location_lat DOUBLE PRECISION;
ALTER TABLE articles ADD COLUMN IF NOT EXISTS location_lng DOUBLE PRECISION;
ALTER TABLE articles ADD COLUMN IF NOT EXISTS global_emotion TEXT;
ALTER TABLE articles ADD COLUMN IF NOT EXISTS counter_argument TEXT;
ALTER TABLE articles ADD COLUMN IF NOT EXISTS causal_links JSONB;

CREATE INDEX IF NOT EXISTS articles_global_emotion_idx ON articles (global_emotion);

-- 003_ads_schema.up.sql
CREATE TABLE IF NOT EXISTS ads (
    id UUID PRIMARY KEY,
    type VARCHAR(50) NOT NULL, -- native, banner
    title TEXT NOT NULL,
    body TEXT,
    media_url TEXT,
    target_url TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    priority INT DEFAULT 5,
    impressions BIGINT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ads_active ON ads(is_active, priority DESC);
