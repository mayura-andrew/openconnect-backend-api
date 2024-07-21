CREATE INDEX IF NOT EXISTS ideas_title_idx ON ideas USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS ideas_genres_idx ON ideas USING GIN (tags);