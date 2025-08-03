-- +goose Up
create table users (
    id UUID primary key default gen_random_uuid(),
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    email TEXT not null unique
);

-- +goose Down
drop table users;
