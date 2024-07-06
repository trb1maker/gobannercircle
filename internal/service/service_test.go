package service_test

import (
	"context"
	"errors"
	"net"
	"strconv"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/trb1maker/gobannercircle/internal/service"
	"github.com/trb1maker/gobannercircle/internal/service/api"
	"github.com/trb1maker/gobannercircle/internal/service/mocks"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestService(t *testing.T) {
	const (
		host = "127.0.0.1"
		port = 8088
	)
	testErr := errors.New("test error")

	app := mocks.NewApp(t)
	app.On("BannerOn", mock.Anything, 1, 1).Return(nil).Times(1)
	app.On("BannerOff", mock.Anything, 10, 11).Return(nil).Times(1)
	app.On("Banner", mock.Anything, 42, 14).Return(34, nil).Times(1)
	app.On("Click", mock.Anything, 1, 1, 1).Return(testErr).Times(1)

	s := service.NewService(app, host, port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan struct{})

	go func() {
		s.Stop(ctx)
	}()

	go func() {
		require.NoError(t, s.Start())
		close(done)
	}()

	conn, err := grpc.NewClient(
		net.JoinHostPort(host, strconv.Itoa(int(port))),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, conn.Close())
	}()

	client := api.NewBannerRotationClient(conn)
	_, err = client.AddBanner(ctx, &api.SlotConfigRequest{
		SlotId:   1,
		BannerId: 1,
	})
	require.NoError(t, err)

	_, err = client.RemoveBanner(ctx, &api.SlotConfigRequest{
		SlotId:   10,
		BannerId: 11,
	})
	require.NoError(t, err)

	res, err := client.Find(ctx, &api.BannerRequest{
		SlotId:  42,
		GroupId: 14,
	})
	require.NoError(t, err)
	require.Equal(t, int64(34), res.GetBannerId())

	_, err = client.Success(ctx, &api.SuccessRequest{
		SlotId:   1,
		BannerId: 1,
		GroupId:  1,
	})
	require.Error(t, err)
	cancel()

	<-done
}
