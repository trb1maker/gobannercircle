-- +goose Up
-- +goose StatementBegin
create table if not exists app.banners (
    banner_id   integer primary key,
    description text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists app.banners;
-- +goose StatementEnd
