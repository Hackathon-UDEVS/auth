CREATE TABLE users (
    id uuid PRIMARY KEY not null ,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('client', 'contractors')), -- Enum-like constraint for role
    created_at timestamp default now() not NULL,
    updated_at timestamp default now() not NULL,
    deleted_at bigint default 0 not NULL
);
