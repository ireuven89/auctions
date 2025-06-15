-- +goose Up

create table if not exists auctions
(
    id              varchar(36) primary key,
    description     varchar(255) not null,
    regions         json         not null,
    initial_offer   integer unsigned,
    current_bid     integer unsigned,
    seller_id       varchar(36)  not null,
    winner_id       varchar(36)  not null default '',
    bidders_count   integer unsigned,
    current_bidders integer unsigned,
    status          varchar(255),
    created_at      timestamp             default current_timestamp,
    updated_at      timestamp             default current_timestamp
);

create table items
(
    id          varchar(36) primary key,
    description text,
    auction_id  varchar(36),
    created_at  timestamp default current_timestamp,
    updated_at  timestamp default current_timestamp,
    foreign key (auction_id) references auctions (id)
);

create table item_pictures
(
    id            varchar(36) primary key,
    download_link text,
    item_id       varchar(36),
    created_at    timestamp default current_timestamp,
    updated_at    timestamp default current_timestamp,
    foreign key (item_id) references items (id)
);

create table if not exists bidders
(
    id         varchar(36) primary key,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);

create table bid
(
    id         varchar(36) primary key,
    auction_id varchar(36),
    bidder_id  varchar(36),
    bid        bigint unsigned default 0,
    foreign key (auction_id) references auctions (id),
    foreign key (bidder_id) references bidders (id)
);

