create table users (
	id uuid primary key,
	created_at time not null default 0,
	name text not null unique,
	password_hash blob not null,
	permissions user_permissions
);

create table user_sessions (
	id uuid primary key,
	created_at time not null default 0,
	expires_at time not null default 0,
	user_id uuid not null,
	secret blob not null,

	foreign key (user_id) references users(id) on update cascade on delete cascade
);