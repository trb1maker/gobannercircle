-- +goose Up
-- +goose StatementBegin
create table if not exists app.user_groups (
    group_id    serial primary key,
    description text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists app.user_groups;
-- +goose StatementEnd
