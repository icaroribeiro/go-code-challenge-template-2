DROP TABLE IF EXISTS auths;

-- Table auths
CREATE TABLE IF NOT EXISTS auths (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT auths_pkey PRIMARY KEY (id),
    CONSTRAINT auths_user_id_key UNIQUE (user_id),
    CONSTRAINT auths_user_id_users_id_foreign FOREIGN KEY (user_id)
        REFERENCES users (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
);