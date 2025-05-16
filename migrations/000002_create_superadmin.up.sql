INSERT INTO users (
    id,
    phone_number,
    password,
    first_name,
    last_name,
    email,
    status,
    role,
    created_at,
    updated_at
) VALUES (
    gen_random_uuid(),
    '09000000000',
    '$2a$10$oJ6KFVfrHQXCekoBz..X4uCkCzz7h7uMWbq8cxbj1G3qWiLw2OId6',
    'Super',
    'Admin',
    'admin@example.com',
    0,
    1,
    NOW(),
    NOW()
); 

