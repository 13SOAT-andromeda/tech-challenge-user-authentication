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
    '$2a$10$/1jabFqWFbWzAfIEIkPOrOVJW8EiBJQv3HRCrHOxKXAgbpxChV9dC', -- password: admin123
    'admin',
    NOW(),
    NOW(),
    NULL
) ON CONFLICT (document) DO NOTHING;
