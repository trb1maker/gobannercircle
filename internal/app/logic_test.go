//go:build logic

package app_test

import (
	"context"
	"database/sql"
	"math/rand"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/trb1maker/gobannercircle/internal/app"
	"github.com/trb1maker/gobannercircle/internal/app/mocks"
	"github.com/trb1maker/gobannercircle/internal/storage/sqlite"
)

const (
	iterCount  = 100_000
	maxSlots   = 10
	maxBanners = 15
	maxGroups  = 20
)

var dbName = "../../db/implementation.sqlite"

func migrations(t *testing.T) {
	t.Helper()

	os.Remove(dbName)

	conn, err := sql.Open("sqlite3", dbName)
	require.NoError(t, err)
	defer conn.Close()

	require.NoError(t, goose.SetDialect("sqlite"))
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	require.NoError(t, goose.UpContext(ctx, conn, "../../migrations/sqlite"))
}

func generator(t *testing.T) {
	conn, err := sql.Open("sqlite3", dbName)
	require.NoError(t, err)
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
	})
	require.NoError(t, err)
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `insert into slots (slot_id) values ($1)`)
	require.NoError(t, err)

	for i := 0; i < maxSlots; i++ {
		_, err := stmt.ExecContext(ctx, i)
		require.NoError(t, err)
	}
	require.NoError(t, stmt.Close())

	stmt, err = tx.PrepareContext(ctx, `insert into banners (banner_id) values ($1)`)
	require.NoError(t, err)

	for i := 0; i < maxBanners; i++ {
		_, err := stmt.ExecContext(ctx, i)
		require.NoError(t, err)
	}
	require.NoError(t, stmt.Close())

	stmt, err = tx.PrepareContext(ctx, `insert into user_groups (group_id) values ($1)`)
	require.NoError(t, err)

	for i := 0; i < maxGroups; i++ {
		_, err := stmt.ExecContext(ctx, i)
		require.NoError(t, err)
	}
	require.NoError(t, stmt.Close())

	_, err = tx.ExecContext(ctx, `
		insert into actions (slot_id, banner_id, group_id)
		select
			slot_id,
			banner_id,
			group_id
		from slots
		join banners
		join user_groups
	`)
	require.NoError(t, err)

	require.NoError(t, tx.Commit())
}

func TestMainLogic(t *testing.T) {
	migrations(t)
	generator(t)

	storage, err := sqlite.NewSQLiteStorage(dbName)
	require.NoError(t, err)
	defer storage.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	require.NoError(t, storage.Connect(ctx))

	ntf := mocks.NewNotifier(t)
	ntf.On("Notify", mock.Anything, mock.Anything).Return(nil)

	app := app.NewApp(storage, ntf)

	randomizer := rand.New(rand.NewSource(time.Now().UnixMilli()))

	var slotID, bannerID, groupID int

	for range iterCount {
		slotID = randomizer.Intn(maxSlots)
		groupID = randomizer.Intn(maxGroups)

		bannerID, err = app.Banner(ctx, slotID, groupID)
		require.NoError(t, err)

		p := randomizer.Float32()
		if p > 0.8 || (p > 0.5 && bannerID == 7 && (groupID == 3 || groupID == 5)) {
			require.NoError(t, app.Click(ctx, slotID, bannerID, groupID))
		}
	}
}
