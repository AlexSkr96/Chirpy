-- name: CreateChirp :one
insert into chirps (user_id, body) values ($1, $2) returning *;

-- name: GetAllChirps :many
select * from chirps order by created_at;

-- name: GetChirpByID :one
select * from chirps where id = $1;
