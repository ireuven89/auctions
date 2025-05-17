-- +goose Up

alter table auctions add index idx_name (name);