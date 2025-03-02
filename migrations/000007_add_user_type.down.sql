ALTER TABLE users
  DROP CONSTRAINT valid_user_type;

ALTER TABLE users
  DROP COLUMN user_type;

ALTER TABLE users
  RENAME COLUMN username TO name;
