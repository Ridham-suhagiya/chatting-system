CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL UNIQUE CHECK (LENGTH(password) >= 6),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);