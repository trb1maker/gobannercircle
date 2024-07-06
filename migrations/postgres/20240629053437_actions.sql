-- +goose Up
-- +goose StatementBegin
create table if not exists actions (
    id        serial  primary key,
    slot_id   integer not null references slots(slot_id) on delete cascade,
    banner_id integer not null references banners(banner_id) on delete cascade,
    group_id  integer not null references user_groups(group_id) on delete cascade,
    views     integer not null default 1,
    clicks    integer not null default 0,
    
    unique(slot_id, banner_id, group_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists actions;
-- +goose StatementEnd
