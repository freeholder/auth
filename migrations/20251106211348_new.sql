-- +goose Up
create table users (
    id serial primary key,
    name varchar(255) not null,
    email varchar(255) not null,
    password varchar(255) not null,
    role integer default 0,
    created_at timestamp not null,
    updated_at timestamp default now()
);
-- +goose Down
drop table users;