//go:build postgres

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
	"github.com/trb1maker/gobannercircle/internal/storage"
)

func Test(t *testing.T) {
	suite.Run(t, new(postgresSuite))
}

type postgresSuite struct {
	suite.Suite
	storage *Storage
	db      *sql.DB
}

func (s *postgresSuite) SetupSuite() {
	const (
		host   = "192.168.0.103"
		port   = 5432
		option = "webhook"
	)

	var err error

	s.db, err = sql.Open("pgx", fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s", host, port, option, option, option))
	s.Require().NoError(err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Require().NoError(goose.SetDialect("postgres"))
	s.Require().NoError(goose.UpContext(ctx, s.db, "../../../migrations/postgres"))

	_, err = s.db.ExecContext(ctx, `
		insert into slots (slot_id)
		values (1), (2), (3);

		insert into banners (banner_id)
		values (1), (2), (3);

		insert into user_groups (group_id)
		values (1), (2), (3);
	`)
	s.Require().NoError(err)

	s.storage, err = NewPostgresStorage(host, port, option, option, option)
	s.Require().NoError(err)

	s.Require().NoError(s.storage.Connect(ctx))
}

func (s *postgresSuite) TearDownSuite() {
	s.Require().NoError(goose.Reset(s.db, "../../../migrations/postgres"))
	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.db.Close())
}

func (s *postgresSuite) TearDownTest() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, `
		delete from actions;
	`)
	s.Require().NoError(err)
}

func (s *postgresSuite) TestBannerOn() {
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

func (s *postgresSuite) TestBannerOff() {
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

func (s *postgresSuite) Count() {
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

	s.Require().NoError(s.storage.CountView(ctx, slotID, bannerID, groupID))

	row := s.db.QueryRowContext(ctx, `
		select
			views
		from actions
		where slot_id = $1 and banner_id = $2 and group_id = $3
	`, slotID, bannerID, groupID)

	var count int
	s.Require().NoError(row.Scan(&count))
	s.Require().Equal(2, count)

	s.Require().NoError(s.storage.CountClick(ctx, slotID, bannerID, groupID))

	row = s.db.QueryRowContext(ctx, `
		select
			clicks
		from actions
		where slot_id = $1 and banner_id = $2 and group_id = $3
	`, slotID, bannerID, groupID)

	s.Require().NoError(row.Scan(&count))
	s.Require().Equal(1, count)
}

func (s *postgresSuite) TestStats() {
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
