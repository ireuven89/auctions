-- +goose Up

alter table bidders add column item varchar(2000) not null;


