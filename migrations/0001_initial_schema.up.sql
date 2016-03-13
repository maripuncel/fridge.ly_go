create table if not exists purchased_items (
	id serial,
	item_id integer not null,
	quantity integer not null,
	unit varchar(50) not null,
	cost decimal (10, 2) not null,
	purchased_at timestamp with time zone default current_timestamp,
	primary key (id)
);

create table if not exists used_items (
	id serial,
	item_id integer not null,
	quantity integer not null,
	unit varchar(50) not null,
	used_at timestamp with time zone default current_timestamp,
	primary key (id)
);

create table if not exists items (
	id serial,
	name varchar(50) not null,
	created_at timestamp with time zone default current_timestamp,
	primary key (id)
)
