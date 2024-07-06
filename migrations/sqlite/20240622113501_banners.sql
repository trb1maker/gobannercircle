-- +goose Up
-- +goose StatementBegin
create table if not exists banners (
    banner_id   integer primary key,
    description text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists banners;
-- +goose StatementEnd
