-- +goose Up

create table if not exists auctions
(
    id  varchar(36) primary key,
    name varchar(255)
);

create table if not exists bidders
(
    id varchar(36),
    name varchar(255)
);