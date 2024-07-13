-- +goose Up
-- +goose StatementBegin
create table if not exists app.slots (
    slot_id     serial primary key,
    description text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists app.slots;
-- +goose StatementEnd
