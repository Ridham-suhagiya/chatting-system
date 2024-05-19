CREATE TABLE chat_links (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    generated_by_user_id UUID NOT NULL,
    link_code VARCHAR(255) NOT NULL,
    expiration_date TIMESTAMP DEFAULT NULL,
    num_users INT DEFAULT 0, 
    FOREIGN KEY (generated_by_user_id) REFERENCES Users(id)
);