-- +goose Up
-- +goose StatementBegin
create table if not exists user_groups (
    group_id    integer primary key,
    description text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists user_groups;
-- +goose StatementEnd
