-- +goose Up
-- +goose StatementBegin
create index if not exists slot_group_idx on actions(slot_id, group_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists slot_group_idx;
-- +goose StatementEnd
