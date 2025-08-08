-- +goose Up
create table refresh_tokens (
    token varchar(255) primary key,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    user_id uuid not null references users(id) on delete cascade,
    expires_at timestamp not null default now() + interval '60 day',
    revoked_at timestamp
);

-- +goose Down
drop table refresh_tokens;
