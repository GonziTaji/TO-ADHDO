/*
drop table tag_of_task;
drop table tasks;
drop table tags;
drop table users;
*/

CREATE TABLE IF NOT EXISTS Users (
    user_id serial PRIMARY KEY,
	username VARCHAR (255) UNIQUE NOT NULL,
	created_on TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS Tags (
    tag_id serial PRIMARY KEY,
    user_id int NOT NULL,
    name VARCHAR (255) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (user_id)
);

CREATE TABLE IF NOT EXISTS Tasks (
    task_id serial PRIMARY KEY,
    user_id int NOT NULL,
    name VARCHAR (255) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (user_id)
);

CREATE TABLE IF NOT EXISTS Tag_of_Task (
    tag_of_task_id serial PRIMARY KEY,
    task_id int NOT NULL,
    tag_id int NOT NULL,
    FOREIGN KEY (task_id) REFERENCES tasks (task_id),
    FOREIGN KEY (tag_id) REFERENCES tags (tag_id)
);

insert into users (username, created_on) values ('gonzalo', now());

insert into tags
(user_id, name)
values 
(1, 'gaming'),
(1, 'pato'),
(1, 'compras'),
(1, 'aseo'),
(1, 'pieza'),
(1, 'trabajo');