CREATE TABLE chat_links (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    link_code VARCHAR(255) NOT NULL,
    expiry_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    num_users INT DEFAULT 0, 
    FOREIGN KEY (user_id) REFERENCES Users(id)
);