CREATE TABLE app_user(
id integer primary key autoincrement,
sn varchar(8) not null unique,
name varchar(6),
email varchar(50),
hashed_password char(60),
created timestamp default (datetime('now','localtime'))
);

CREATE TABLE app_role(
id integer primary key autoincrement,
name varchar(20) NOT NULL unique,
description varchar(50)
);
