-- +goose Up

alter table auctions modify column description text not null;