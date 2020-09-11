create table if not exists trainers
(
	id string not null
		constraint trainers_pk
			primary key
);

create table if not exists pokemon
(
    id integer not null
        constraint pokemon_pk
            primary key autoincrement,
    trainer string not null
        constraint "pokemon_train_trainers.id_fk"
            references trainers ("id"),
    pokemon_id integer not null,
    name text
);

create index pokemon_trainer_index
    on pokemon (trainer);
