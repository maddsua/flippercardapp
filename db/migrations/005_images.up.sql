create table images (
	id text primary key,
	created_at time not null default 0,
	mimetype text not null,
	source_name text not null,
	source_sha512_hash blob not null,
	data blob not null,
	data_size integer not null,
	data_sha512_hash blob not null
);
