-- +goose Up
create table chirps (
    id UUID primary key default gen_random_uuid(),
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    body varchar(140) not null,
    user_id UUID references users(id) on delete cascade
);

-- +goose Down
drop table chirps;
