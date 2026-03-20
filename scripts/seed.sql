INSERT INTO "User" (
    "name",
    "email",
    "contact",
    "document",
    "is_active",
    "password",
    "role",
    "created_at",
    "updated_at",
    "deleted_at"
) VALUES (
    'Admin User',
    'admin@admin.com',
    '+5511999999999',
    '123.456.789-00',
    true,
    '$2b$12$PZLp69vH.c/.GSu9RpWwYe8BPUkGkSF6.7rzAOpMA.PDfIAxzbr0.',
    'admin',
    NOW(),
    NOW(),
    NULL
) ON CONFLICT (document) DO NOTHING;
