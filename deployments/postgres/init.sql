create table if not exists users
(
    id uuid primary key default gen_random_uuid(),
    first_name text,
    last_name text,
    nickname text,
    password text,
    email text,
    country text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);
