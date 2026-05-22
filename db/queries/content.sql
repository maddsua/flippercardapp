-- name: GetDeckById :one
select * from decks
where id = sqlc.arg(id);

-- name: GetDecksBatch :many
select
	distinct decks.*,
	count(cards.id) as size
from decks
	inner join cards on cards.deck_id = decks.id
where (sqlc.narg(ids_set) is null or decks.id in (
	select value from json_each(sqlc.narg(ids_set))
)) and (decks.collection_id = sqlc.narg(collection_id)
	or sqlc.narg(collection_id) is null
) and (sqlc.narg(visibility_set) is null or decks.visibility in (
	select value from json_each(sqlc.narg(visibility_set))
))
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
where (sqlc.narg(ids_set) is null or collections.id in (
	select value from json_each(sqlc.narg(ids_set))
)) and (sqlc.narg(visibility_set) is null or collections.visibility in (
	select value from json_each(sqlc.narg(visibility_set))
))
group by collections.id
order by collections.created_at desc
limit sqlc.arg(limit) offset sqlc.arg(offset);

-- name: GetCollectionSearchBatch :many
select id, name from collections
where (sqlc.narg(visibility_set) is null or collections.visibility in (
	select value from json_each(sqlc.narg(visibility_set))
))
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
	description,
	visibility
) values (
	sqlc.arg(id),
	sqlc.arg(created_at),
	sqlc.arg(updated_at),
	sqlc.arg(name),
	sqlc.narg(description),
	sqlc.arg(visibility)
) returning *;

-- name: UpdateCollection :one
update collections
set
	updated_at = sqlc.arg(updated_at),
	name = sqlc.arg(name),
	description = sqlc.arg(description),
	visibility = sqlc.arg(visibility)
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
	description,
	visibility
) values (
	sqlc.arg(id),
	sqlc.arg(collection_id),
	sqlc.arg(created_at),
	sqlc.arg(updated_at),
	sqlc.arg(name),
	sqlc.narg(description),
	sqlc.arg(visibility)
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

-- name: SetDeckUpdateTime :one
update decks
set updated_at = sqlc.arg(updated_at)
where id = sqlc.arg(id)
returning *;

-- name: UpdateDeckMetadata :one
update decks
set
	collection_id = coalesce(sqlc.narg(collection_id), collection_id),
	name = sqlc.arg(name),
	description = sqlc.arg(description),
	visibility = sqlc.arg(visibility)
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
