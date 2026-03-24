-- Seed users for each role
-- Password for all: Admin123!
-- Hash: $2a$10$xV39RrwZ/TiSjR08EvmkxuKLqSopxJLe9HrwmA7HaYsd1VorDMlm2

-- Create table if GORM AutoMigrate hasn't run yet
CREATE TABLE IF NOT EXISTS "User" (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name       TEXT NOT NULL DEFAULT '',
    email      TEXT NOT NULL DEFAULT '',
    contact    TEXT NOT NULL DEFAULT '',
    document   TEXT NOT NULL,
    is_active  BOOLEAN DEFAULT true,
    password   TEXT NOT NULL DEFAULT '',
    role       TEXT NOT NULL DEFAULT '',
    street     TEXT DEFAULT '',
    number     TEXT DEFAULT '',
    complement TEXT DEFAULT '',
    city       TEXT DEFAULT '',
    state      TEXT DEFAULT '',
    zip_code   TEXT DEFAULT ''
);
CREATE UNIQUE INDEX IF NOT EXISTS uni_user_document ON "User" (document);

INSERT INTO "User" (
    "name", "email", "contact", "document",
    "is_active", "password", "role",
    "street", "number", "complement", "city", "state", "zip_code",
    "created_at", "updated_at", "deleted_at"
) VALUES
(
    'Customer User',
    'customer@example.com',
    '+5511911111111',
    '11122233344',
    true,
    '$2a$10$xV39RrwZ/TiSjR08EvmkxuKLqSopxJLe9HrwmA7HaYsd1VorDMlm2',
    'customer',
    '', '', '', '', '', '',
    NOW(), NOW(), NULL
),
(
    'Attendant User',
    'attendant@example.com',
    '+5511922222222',
    '22233344455',
    true,
    '$2a$10$xV39RrwZ/TiSjR08EvmkxuKLqSopxJLe9HrwmA7HaYsd1VorDMlm2',
    'attendant',
    '', '', '', '', '', '',
    NOW(), NOW(), NULL
),
(
    'Mechanic User',
    'mechanic@example.com',
    '+5511933333333',
    '33344455566',
    true,
    '$2a$10$xV39RrwZ/TiSjR08EvmkxuKLqSopxJLe9HrwmA7HaYsd1VorDMlm2',
    'mechanic',
    '', '', '', '', '', '',
    NOW(), NOW(), NULL
),
(
    'Administrator User',
    'administrator@example.com',
    '+5511944444444',
    '44455566677',
    true,
    '$2a$10$xV39RrwZ/TiSjR08EvmkxuKLqSopxJLe9HrwmA7HaYsd1VorDMlm2',
    'administrator',
    '', '', '', '', '', '',
    NOW(), NOW(), NULL
)
ON CONFLICT (document) DO UPDATE SET
    password   = EXCLUDED.password,
    name       = EXCLUDED.name,
    email      = EXCLUDED.email,
    is_active  = EXCLUDED.is_active,
    role       = EXCLUDED.role,
    updated_at = NOW();
