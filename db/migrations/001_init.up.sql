create table collections (
	id uuid primary key,
	created_at time not null default 0,
	updated_at time not null default 0,
	name text not null,
	description text
);

create table decks (
	id uuid primary key,
	collection_id uuid not null,
	created_at time not null default 0,
	updated_at time not null default 0,
	name text not null,
	description text,

	foreign key (collection_id) references collections(id) on update cascade on delete cascade
);

create table cards (
	id uuid primary key,
	deck_id uuid not null,
	created_at time not null default 0,
	updated_at time not null default 0,
	content card_node_content not null,

	foreign key (deck_id) references decks(id) on update cascade on delete cascade
);
