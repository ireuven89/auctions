-- +goose Up

alter table auctions add constraint unique key unique_name(name);

alter table bidders add constraint unique key  unique_name(name);