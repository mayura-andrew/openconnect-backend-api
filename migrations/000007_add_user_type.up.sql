ALTER TABLE users
    RENAME COLUMN name TO user_name;

ALTER TABLE users
    ADD COLUMN user_type TEXT NOT NULL DEFAULT 'normal';

ALTER TABLE users
    ADD CONSTRAINT valid_user_type CHECK (user_type IN ('normal', 'admin', 'google'));
