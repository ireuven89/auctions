-- +goose Up

alter table bidders add column download_link varchar(2000) not null ;