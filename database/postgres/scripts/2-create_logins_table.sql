DROP TABLE IF EXISTS logins;

-- Table logins
CREATE TABLE IF NOT EXISTS logins (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    username text NOT NULL,
    password text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT logins_pkey PRIMARY KEY (id),
    CONSTRAINT logins_user_id_key UNIQUE (user_id),
    CONSTRAINT logins_user_id_users_id_foreign FOREIGN KEY (user_id)
        REFERENCES users (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
);