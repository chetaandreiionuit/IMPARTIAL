-- Up Migration

CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE articles (
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

-- Index for fast cosine similarity search
CREATE INDEX ON articles USING hnsw (embedding vector_cosine_ops);
