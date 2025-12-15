-- Create Ads Table

CREATE TABLE ads (
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

CREATE INDEX idx_ads_active ON ads(is_active, priority DESC);
