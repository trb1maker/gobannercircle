-- +goose Up
-- +goose StatementBegin
create schema if not exists app;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop schema if exists app;
-- +goose StatementEnd
