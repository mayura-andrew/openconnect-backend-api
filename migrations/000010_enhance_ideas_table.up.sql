-- Add new columns to ideas table
ALTER TABLE ideas 
    ADD COLUMN IF NOT EXISTS learning_outcome TEXT,
    ADD COLUMN IF NOT EXISTS recommended_level VARCHAR(50),
    ADD COLUMN IF NOT EXISTS github_link TEXT,
    ADD COLUMN IF NOT EXISTS website_link TEXT,
    ADD COLUMN IF NOT EXISTS user_id UUID;

-- Add the foreign key constraint for user_id (after adding the column)
ALTER TABLE ideas 
    ADD CONSTRAINT fk_ideas_user 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Add PDF column if it doesn't exist
ALTER TABLE ideas 
    ADD COLUMN IF NOT EXISTS pdf TEXT;

-- Create indexes for faster searches
CREATE INDEX IF NOT EXISTS idx_ideas_user_id ON ideas(user_id);
CREATE INDEX IF NOT EXISTS idx_ideas_status ON ideas(status);
CREATE INDEX IF NOT EXISTS idx_ideas_category ON ideas(category);
CREATE INDEX IF NOT EXISTS idx_ideas_tags ON ideas USING gin(tags);