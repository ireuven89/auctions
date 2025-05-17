-- +goose Up

create table bidders(
    id varchar(36) primary key,
    name varchar(100) unique key
);

create index idx_name on bidders(name);