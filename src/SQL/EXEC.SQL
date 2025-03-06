CREATE TABLE IF NOT EXISTS PRJ
(
    id      integer primary key autoincrement,
    id_item integer,
    name    varchar not null,
    comment varchar,
    tags 	varchar
);


CREATE TABLE IF NOT EXISTS OBJ
(
    id          integer primary key autoincrement,
    id_prj	    integer,
    id_parent   integer,
    name        varchar not null,
    comment     varchar,
    object_type integer not null, -- 0 - folder / 1 - file
    tags 		varchar,
    exec        varchar,
    output      varchar,
    FOREIGN KEY (id_prj) REFERENCES PRJ (id) ON DELETE CASCADE,
    FOREIGN KEY (id_parent) REFERENCES OBJ (id) ON DELETE CASCADE
    );


CREATE TABLE IF NOT EXISTS SRC
(
    id      	integer primary key autoincrement,
    id_prj	    integer,
    id_file		integer NOT null,
    num         integer,
    line    	varchar,
    comment 	varchar,
    tags 		varchar,
    todo        boolean default False,
    FOREIGN KEY (id_file) REFERENCES OBJ (id) ON DELETE CASCADE
    );