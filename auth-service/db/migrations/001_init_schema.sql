-- +goose Up

create table users(
    id varchar(36) primary key,
    name varchar(36),
    password varchar(60),
    email varchar(255)
);