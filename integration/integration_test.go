//go:build integration

package integration_test

import (
	"context"
	"net"
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/require"
	"github.com/trb1maker/gobannercircle/internal/service/api"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	iterCount  = 100_000
	maxSlots   = 10
	maxBanners = 15
	maxGroups  = 20
)

func generator(t *testing.T) {
	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "app",
		User:     "app",
		Password: "app",
	})
	require.NoError(t, err)
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tx, err := conn.BeginEx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()

	for i := 0; i < maxSlots; i++ {
		_, err := tx.ExecEx(ctx, `insert into slots (slot_id) values ($1)`, nil, i)
		require.NoError(t, err)
	}

	for i := 0; i < maxBanners; i++ {
		_, err := tx.ExecEx(ctx, `insert into banners (banner_id) values ($1)`, nil, i)
		require.NoError(t, err)
	}

	for i := 0; i < maxGroups; i++ {
		_, err := tx.ExecEx(ctx, `insert into user_groups (group_id) values ($1)`, nil, i)
		require.NoError(t, err)
	}

	require.NoError(t, tx.Commit())
}

func TestIntegration(t *testing.T) {
	generator(t)

	conn, err := grpc.NewClient(
		net.JoinHostPort("localhost", "8088"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := api.NewBannerRotationClient(conn)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Interrupt)
	defer cancel()

	for slotID := 0; slotID < maxSlots; slotID++ {
		for bannerID := 0; bannerID < maxBanners; bannerID++ {
			_, err = client.AddBanner(ctx, &api.SlotConfigRequest{
				SlotId:   int64(slotID),
				BannerId: int64(bannerID),
			})
			require.NoError(t, err)
		}
	}

	randomizer := rand.New(rand.NewSource(uint64(time.Now().Second())))

	var slotID, bannerID, groupID int
	for range iterCount {
		slotID = randomizer.Intn(maxSlots)
		groupID = randomizer.Intn(maxGroups)

		banner, err := client.Find(ctx, &api.BannerRequest{
			SlotId:  int64(slotID),
			GroupId: int64(groupID),
		})
		require.NoError(t, err)
		bannerID = int(banner.GetBannerId())

		p := randomizer.Float32()
		if p > 0.8 || (p > 0.5 && bannerID == 7 && (groupID == 3 || groupID == 5)) {
			_, err = client.Success(ctx, &api.SuccessRequest{
				SlotId:   int64(slotID),
				BannerId: int64(bannerID),
				GroupId:  int64(groupID),
			})

			require.NoError(t, err)
		}
	}
}
