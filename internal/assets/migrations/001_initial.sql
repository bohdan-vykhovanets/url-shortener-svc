-- +migrate Up

create table shortened_urls
(
    id  bigserial primary key not null,
    code text unique not null,
    long_url text not null,
    created_at timestamp without time zone
);

create index shortened_urls_index on shortened_urls (code);

-- +migrate Down

drop table shortened_urls;
