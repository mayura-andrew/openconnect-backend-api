






CREATE TABLE IF NOT EXISTS ideas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP(0) with TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) with TIME ZONE NOT NULL DEFAULT NOW(),
    submitted_by UUID NOT NULL,
    idea_source_id UUID NOT NULL,
    category VARCHAR(100) NOT NULL,
    tags TEXT NOT NULL,
    upvotes INT NOT NULL DEFAULT 0,
    downvotes INT NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    version INT NOT NULL DEFAULT 1,
    comments TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    interested_users UUID[] NOT NULL DEFAULT ARRAY[]::UUID[] -- Adjusted for UUID
);