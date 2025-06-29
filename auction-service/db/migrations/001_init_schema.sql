-- +goose Up

create table if not exists auctions
(
    id  varchar(36) primary key,
    name varchar(255),
    description text,
    user_id varchar(36),
    active boolean default false,
    end_time timestamp,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);