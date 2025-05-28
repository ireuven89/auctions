-- +goose Up
alter table users modify column name varchar(600) not null;
