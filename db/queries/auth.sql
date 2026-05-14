-- name: UserCount :one
select count(1) from users;

-- name: InsertUser :one
insert into users (
	id,
	created_at,
	name,
	password_hash,
	permissions
) values (
	sqlc.arg(id),
	sqlc.arg(created_at),
	sqlc.arg(name),
	sqlc.arg(password_hash),
	sqlc.arg(permissions)
) returning *;

-- name: GetUserByID :one
select * from users
where id = sqlc.arg(id);

-- name: GetUserByName :one
select * from users
where name = sqlc.arg(name);

-- name: InsertSession :one
insert into user_sessions (
	id,
	created_at,
	expires_at,
	user_id,
	secret
) values (
	sqlc.arg(id),
	sqlc.arg(created_at),
	sqlc.arg(expires_at),
	sqlc.arg(user_id),
	sqlc.arg(secret)
) returning *;

-- name: GetSession :one
select * from user_sessions
where id = sqlc.arg(id);

-- name: InvalidateSession :exec
update user_sessions
set expires_at = 0, secret = null
where id = sqlc.arg(id);
