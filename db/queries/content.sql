-- name: GetDeckById :one
select * from decks
where id = sqlc.arg(id);

-- name: GetDecksBatch :many
select
	distinct decks.*,
	count(distinct cards.id) as size
from decks
	left join cards on cards.deck_id = decks.id
where (decks.id = sqlc.narg(id) or sqlc.narg(id) is null)
	and (decks.collection_id = sqlc.narg(collection_id) or sqlc.narg(collection_id) is null)
limit sqlc.arg(limit) offset sqlc.arg(offset);

-- name: GetDeckCards :many
select * from cards
where deck_id = sqlc.arg(deck_id);

-- name: GetCollectionById :one
select * from collections
where id = sqlc.arg(id);

-- name: GetCollectionBatch :many
select
	collections.*,
	count(distinct decks.id) as size
from collections
	left join decks on decks.collection_id = collections.id
where collections.id = sqlc.narg(id) or sqlc.narg(id) is null
limit sqlc.arg(limit) offset sqlc.arg(offset);

-- name: InsertCollection :one
insert into collections (
	id,
	created_at,
	updated_at,
	name,
	description
) values (
	sqlc.arg(id),
	sqlc.arg(created_at),
	sqlc.arg(updated_at),
	sqlc.arg(name),
	sqlc.narg(description)
) returning *;
