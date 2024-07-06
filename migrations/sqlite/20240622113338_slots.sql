-- +goose Up
-- +goose StatementBegin
create table if not exists slots (
    slot_id     integer primary key,
    description text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists slots;
-- +goose StatementEnd
