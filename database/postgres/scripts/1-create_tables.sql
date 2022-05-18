DROP TABLE IF EXISTS auths;
DROP TABLE IF EXISTS logins;
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