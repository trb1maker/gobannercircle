package sqlite

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // lint:revive
	"github.com/trb1maker/gobannercircle/internal/storage"
)

// Storage на базе SQLite предназначен для разработки и тестирования
// работы алгоритма ротации баннеров.
type Storage struct {
	db *sqlx.DB
}

func NewSQLiteStorage(pathToDB string) (*Storage, error) {
	db, err := sqlx.Open("sqlite3", pathToDB)
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *Storage) BannerOn(ctx context.Context, slotID, bannerID int) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		insert or ignore into actions (slot_id, banner_id, group_id)
		select
			$1 slot_id,
			$2 banner_id,
			group_id 
		from user_groups
	`, slotID, bannerID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) BannerOff(ctx context.Context, slotID, bannerID int) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		delete from actions
		where slot_id = $1 and banner_id = $2
	`, slotID, bannerID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) Stats(ctx context.Context, slotID, groupID int) (storage.Stats, error) {
	var stats storage.Stats

	rows, err := s.db.QueryContext(ctx, `
		select
			banner_id,
			views,
			clicks
		from actions
		where slot_id = $1 and group_id = $2
	`, slotID, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var id, views, clicks int
	for rows.Next() {
		if err := rows.Scan(&id, &views, &clicks); err != nil {
			return nil, err
		}
		stats = append(stats, storage.Stat{
			ID:     id,
			Views:  views,
			Clicks: clicks,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

func (s *Storage) Banner(ctx context.Context, slotID, groupID int) (int, error) {
	row := s.db.QueryRowxContext(ctx, `
		select
			banner_id,
			avg(clicks * 1.0 / views) over() + sqrt(2 * ln(sum(views) over()) / views) weight
		from actions
		where slot_id = $1 and group_id = $2
		order by weight desc
		limit 1
	`, slotID, groupID)

	var (
		bannerID int
		weight   float64
	)
	if err := row.Scan(&bannerID, &weight); err != nil {
		return 0, err
	}

	return bannerID, nil
}

func (s *Storage) CountView(ctx context.Context, slotID, bannerID, groupID int) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		update actions set views = views + 1
		where slot_id = $1 and banner_id = $2 and group_id = $3
	`, slotID, bannerID, groupID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) CountClick(ctx context.Context, slotID, bannerID, groupID int) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		update actions set clicks = clicks + 1
		where slot_id = $1 and banner_id = $2 and group_id = $3
	`, slotID, bannerID, groupID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) Close() error {
	return s.db.Close()
}
