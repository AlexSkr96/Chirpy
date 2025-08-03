-- name: CreateChirp :one
insert into chirps (user_id, body) values ($1, $2) returning *;
