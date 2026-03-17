CREATE DATABASE garagedb;

\c garagedb;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    document VARCHAR(14) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert some test data
INSERT INTO users (document, is_active) VALUES ('123.456.789-00', true);
INSERT INTO users (document, is_active) VALUES ('111.222.333-44', false);
