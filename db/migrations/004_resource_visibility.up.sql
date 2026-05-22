alter table collections
	add column visibility resource_visibility not null default 2;

alter table decks
	add column visibility resource_visibility not null default 2;
