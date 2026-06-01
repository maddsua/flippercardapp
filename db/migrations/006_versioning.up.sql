
create table deck_versions (
	id uuid primary key,
	created_at time not null default 0,
	deck_id uuid not null,
	card_count integer not null default 0,
	content deck_version_content not null,
	label text,

	foreign key (deck_id) references decks(id) on update cascade on delete cascade
);

insert into deck_versions (id, deck_id, card_count, content)
select id, id as deck_id, (
	select count(1)
	from cards where cards.deck_id = decks.id
), json_object(
	'cards', (
		select json_group_array(
			json_patch(
				json_object('id', cards.id),
				json(
					cast(
						coalesce(cards.content, '{}') as text
					)
				)
			)
		)
		from cards
		where cards.deck_id = decks.id
	)
) as content from decks;

create index deck_versions_index on deck_versions(deck_id, created_at desc);

/*
	WARNING: This doesn't enforce a foreign key constraint
	as that would be too much hussle to get with not that old versions of SQLite,
	therefore all a fallback to created_at must be used to determine the latest version
	just in case the foreign key gets broken
*/
alter table decks add column latest_version_id uuid;

update decks
set latest_version_id = (
	select id from deck_versions
	where deck_versions.deck_id = decks.id
	order by created_at desc
	limit 1
);

drop table cards;
