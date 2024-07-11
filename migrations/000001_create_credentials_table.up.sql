CREATE TABLE IF NOT EXISTS credentials (
    email VARCHAR(100),
    password VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    CONSTRAINT email_pkey PRIMARY KEY (email)
);