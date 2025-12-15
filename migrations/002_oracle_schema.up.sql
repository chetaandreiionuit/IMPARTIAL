-- Add Causal and Geo columns

ALTER TABLE articles ADD COLUMN IF NOT EXISTS location_lat DOUBLE PRECISION;
ALTER TABLE articles ADD COLUMN IF NOT EXISTS location_lng DOUBLE PRECISION;
ALTER TABLE articles ADD COLUMN IF NOT EXISTS global_emotion TEXT;
ALTER TABLE articles ADD COLUMN IF NOT EXISTS counter_argument TEXT;
ALTER TABLE articles ADD COLUMN IF NOT EXISTS causal_links JSONB; -- Store []CausalLink as JSONB

CREATE INDEX ON articles (global_emotion);
