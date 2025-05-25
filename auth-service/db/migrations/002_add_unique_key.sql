-- +goose Up

alter table users add constraint unique_name unique (name);

alter table users add constraint  unique_email unique (email);