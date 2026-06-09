-- name: GetDeckById :one
select * from decks
where id = sqlc.arg(id);

-- name: GetDecksBatch :many
select
	distinct decks.*,
	deck_versions.card_count as size
from decks
	left join deck_versions on deck_versions.id = decks.latest_version_id
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

-- name: GetDeckVersion :one
select * from deck_versions
where id = sqlc.arg(version_id) and (deck_id = sqlc.narg(deck_id) or sqlc.narg(deck_id) is null);

-- name: GetDeckVersionsBatch :many
select
	id,
	created_at,
	deck_id,
	card_count,
	label
from deck_versions
where deck_id = sqlc.arg(deck_id)
order by created_at desc
limit sqlc.arg(limit) offset sqlc.arg(offset);

-- name: GetDeckLatestVersion :one
select * from deck_versions
where deck_id = sqlc.arg(deck_id)
order by created_at desc
limit 1;

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

-- name: UpdateCollectionChildrenVisibility :execrows
update decks
set visibility = sqlc.arg(new_visibility)
where visibility = sqlc.arg(old_visibility)
	and collection_id = sqlc.arg(collection_id);

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

-- name: InsertDeckVersion :one
insert into deck_versions (
	id,
	created_at,
	deck_id,
	card_count,
	content,
	label
) values (
	sqlc.arg(id),
	sqlc.arg(created_at),
	sqlc.arg(deck_id),
	sqlc.arg(card_count),
	sqlc.arg(content),
	sqlc.arg(label)
) returning *;

-- name: SetDeckLatestVersion :one
update decks
set latest_version_id = sqlc.arg(latest_version_id),
	updated_at = coalesce(sqlc.narg(updated_at), updated_at)
where id = sqlc.arg(deck_id)
returning *;

-- name: UpdateDeckMetadata :one
update decks
set
	collection_id = coalesce(sqlc.narg(collection_id), collection_id),
	name = sqlc.arg(name),
	description = sqlc.arg(description),
	visibility = sqlc.arg(visibility),
	updated_at = coalesce(sqlc.narg(updated_at), updated_at)
where id = sqlc.arg(id)
returning *;

-- name: DeleteDeck :execrows
delete from decks
where id = sqlc.arg(id);

-- name: InsertImage :one
insert into images (
	id,
	created_at,
	mimetype,
	source_name,
	source_sha512_hash,
	data,
	data_size,
	data_sha512_hash
) values (
	sqlc.arg(id),
	sqlc.arg(created_at),
	sqlc.arg(mimetype),
	sqlc.arg(source_name),
	sqlc.arg(source_sha512_hash),
	sqlc.arg(data),
	sqlc.arg(data_size),
	sqlc.arg(data_sha512_hash)
) returning *;

-- name: GetImageByHash :one
select * from images
where source_sha512_hash = sqlc.arg(sha512_hash)
	or data_sha512_hash = sqlc.arg(sha512_hash);

-- name: GetImageById :one
select * from images
where id = sqlc.arg(id);
