package sqlite

import (
	"context"
	"database/sql"
	"os"
	"path"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
	"github.com/trb1maker/gobannercircle/internal/storage"
)

func Test(t *testing.T) {
	suite.Run(t, new(sqliteSuite))
}

type sqliteSuite struct {
	suite.Suite
	storage *Storage
	db      *sql.DB
	path    string
}

func (s *sqliteSuite) SetupSuite() {
	var err error

	s.path = path.Join(os.TempDir(), "storage.db")

	s.db, err = sql.Open("sqlite3", s.path)
	s.Require().NoError(err)
	s.Require().NoError(goose.SetDialect("sqlite"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Require().NoError(goose.UpContext(ctx, s.db, "../../../migrations/sqlite"))

	_, err = s.db.ExecContext(ctx, `
		insert into slots (slot_id)
		values (1), (2), (3);

		insert into banners (banner_id)
		values (1), (2), (3);

		insert into user_groups (group_id)
		values (1), (2), (3);
	`)
	s.Require().NoError(err)

	s.storage, err = NewSQLiteStorage(s.path)
	s.Require().NoError(err)

	s.Require().NoError(s.storage.Connect(ctx))
}

func (s *sqliteSuite) TearDownSuite() {
	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.db.Close())
	s.Require().NoError(os.Remove(s.path))
}

func (s *sqliteSuite) TearDownTest() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, `
		delete from actions;
	`)
	s.Require().NoError(err)
}

func (s *sqliteSuite) TestBannerOn() {
	const (
		slotID   = 2
		bannerID = 3
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Require().NoError(s.storage.BannerOn(ctx, slotID, bannerID))

	row := s.db.QueryRowContext(ctx, `
		select
			count(*)
		from actions
		where slot_id = $1 and banner_id = $2
	`, slotID, bannerID)

	var count int
	s.Require().NoError(row.Scan(&count))
	s.Require().Equal(3, count)
}

func (s *sqliteSuite) TestBannerOff() {
	const (
		slotID   = 2
		bannerID = 3
		groupID  = 2
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, `
		insert into actions (slot_id, banner_id, group_id)
		values ($1, $2, $3)
	`, slotID, bannerID, groupID)
	s.Require().NoError(err)

	s.Require().NoError(s.storage.BannerOff(ctx, slotID, bannerID))

	count := 1
	row := s.db.QueryRowContext(ctx, `
		select
			count(*)
		from actions
		where slot_id = $1 and banner_id = $2
	`, slotID, bannerID)

	s.Require().NoError(row.Scan(&count))
	s.Require().Equal(0, count)
}

func (s *sqliteSuite) Count() {
	const (
		slotID   = 2
		bannerID = 3
		groupID  = 2
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, `
		insert into actions (slot_id, banner_id, group_id)
		values ($1, $2, $3)
	`, slotID, bannerID)
	s.Require().NoError(err)

	s.Require().NoError(s.storage.IncViewCount(ctx, slotID, bannerID, groupID))

	row := s.db.QueryRowContext(ctx, `
		select
			views
		from actions
		where slot_id = $1 and banner_id = $2 and group_id = $3
	`, slotID, bannerID, groupID)

	var count int
	s.Require().NoError(row.Scan(&count))
	s.Require().Equal(2, count)

	s.Require().NoError(s.storage.IncClickCount(ctx, slotID, bannerID, groupID))

	row = s.db.QueryRowContext(ctx, `
		select
			clicks
		from actions
		where slot_id = $1 and banner_id = $2 and group_id = $3
	`, slotID, bannerID, groupID)

	s.Require().NoError(row.Scan(&count))
	s.Require().Equal(1, count)
}

func (s *sqliteSuite) TestStats() {
	const (
		slotID   = 2
		bannerID = 3
		groupID  = 2
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, `
		insert into actions (slot_id, banner_id, group_id)
		values ($1, $2, $3)
	`, slotID, bannerID, groupID)
	s.Require().NoError(err)

	stats, err := s.storage.Stats(ctx, slotID, groupID)
	s.Require().NoError(err)

	s.Require().Len(stats, 1)
	s.Require().Equal(storage.Stat{
		ID:     bannerID,
		Views:  1,
		Clicks: 0,
	}, stats[0])
}
