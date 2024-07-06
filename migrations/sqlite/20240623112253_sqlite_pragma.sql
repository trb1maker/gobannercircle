-- +goose Up
-- +goose StatementBegin
pragma foreign_keys = on;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
pragma foreign_keys = off;
-- +goose StatementEnd
