-- +goose Up

alter table auctions add column bidder_id varchar(36) not null default '';
