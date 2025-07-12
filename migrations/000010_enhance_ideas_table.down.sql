-- Remove indices
DROP INDEX IF EXISTS idx_ideas_tags;
DROP INDEX IF EXISTS idx_ideas_category;
DROP INDEX IF EXISTS idx_ideas_status;
DROP INDEX IF EXISTS idx_ideas_user_id;

-- Remove foreign key constraint
ALTER TABLE ideas DROP CONSTRAINT IF EXISTS fk_ideas_user;

-- Drop the new columns
ALTER TABLE ideas
    DROP COLUMN IF EXISTS learning_outcome,
    DROP COLUMN IF EXISTS recommended_level,
    DROP COLUMN IF EXISTS github_link,
    DROP COLUMN IF EXISTS website_link,
    DROP COLUMN IF EXISTS user_id,
    DROP COLUMN IF EXISTS pdf;