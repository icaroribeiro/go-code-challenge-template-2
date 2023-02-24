DROP TABLE IF EXISTS users;

-- Table users
CREATE TABLE IF NOT EXISTS users (
    id uuid NOT NULL,
    username text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_username_key UNIQUE (username)
);