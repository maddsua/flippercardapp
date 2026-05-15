-- name: GetDeckById :one
select * from decks
where id = sqlc.arg(id);

-- name: GetDecksBatch :many
select
	distinct decks.*,
	count(cards.id) as size
from decks
	inner join cards on cards.deck_id = decks.id
where (decks.id = sqlc.narg(id) or sqlc.narg(id) is null)
	and (decks.collection_id = sqlc.narg(collection_id) or sqlc.narg(collection_id) is null)
group by decks.id
order by decks.created_at desc
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
	count(decks.id) as size
from collections
	left join decks on decks.collection_id = collections.id
where (collections.id = sqlc.narg(id) or sqlc.narg(id) is null)
group by collections.id
order by collections.created_at desc
limit sqlc.arg(limit) offset sqlc.arg(offset);

-- name: GetCollectionSearchBatch :many
select id, name from collections
limit sqlc.arg(limit) offset sqlc.arg(offset);

-- name: CollectionIDExists :one
select exists (
	select 1 from collections
	where id = sqlc.arg(id)
);

-- name: CollectionNameExists :one
select exists (
	select 1 from collections
	where name = sqlc.arg(name)
);

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

-- name: UpdateCollection :one
update collections
set
	updated_at = sqlc.arg(updated_at),
	name = sqlc.arg(name),
	description = sqlc.arg(description)
where id = sqlc.arg(id)
returning *;

-- name: CollectionSize :one
select count(decks.id)
from collections
	left join decks on decks.collection_id = collections.id
where collections.id = sqlc.arg(id);

-- name: DeleteCollection :execrows
delete from collections
where id = sqlc.arg(id);

-- name: InsertDeck :one
insert into decks (
	id,
	collection_id,
	created_at,
	updated_at,
	name,
	description
) values (
	sqlc.arg(id),
	sqlc.arg(collection_id),
	sqlc.arg(created_at),
	sqlc.arg(updated_at),
	sqlc.arg(name),
	sqlc.narg(description)
) returning *;

-- name: InsertCard :exec
insert into cards (
	id,
	deck_id,
	created_at,
	updated_at,
	content
) values (
	sqlc.arg(id),
	sqlc.arg(deck_id),
	sqlc.arg(created_at),
	sqlc.arg(updated_at),
	sqlc.arg(content)
);

-- name: UpdateDeckMetadata :one
update decks
set
	updated_at = sqlc.arg(updated_at),
	collection_id = coalesce(sqlc.narg(collection_id), collection_id),
	name = sqlc.arg(name),
	description = sqlc.arg(description)
where id = sqlc.arg(id)
returning *;

-- name: DeckCardSet :many
select id from cards
where deck_id = sqlc.arg(deck_id);

-- name: UpdateCardContent :execrows
update cards
set
	updated_at = sqlc.arg(updated_at),
	content = sqlc.arg(content)
where id = sqlc.arg(id)
	and deck_id = sqlc.arg(deck_id);

-- name: DeleteCard :exec
delete from cards
where id = sqlc.arg(id);

-- name: DeleteDeck :execrows
delete from decks
where id = sqlc.arg(id);
