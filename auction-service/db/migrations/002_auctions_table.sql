-- +goose Up

alter table auctions add column category varchar(255);