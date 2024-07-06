//go:generate mockery --name Storage
//go:generate mockery --name Notifier
package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/trb1maker/gobannercircle/internal/app"
	"github.com/trb1maker/gobannercircle/internal/app/mocks"
	"github.com/trb1maker/gobannercircle/internal/storage"
)

func TestBannerOn(t *testing.T) {
	stop := errors.New("ok")

	store := mocks.NewStorage(t)
	store.On("BannerOn", mock.Anything, 1, 1).Return(nil)
	store.On("BannerOn", mock.Anything, 1, 2).Return(stop)

	app := app.NewApp(store, nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	require.NoError(t, app.BannerOn(ctx, 1, 1))
	require.ErrorIs(t, app.BannerOn(ctx, 1, 2), stop)
}

func TestBannerOff(t *testing.T) {
	store := mocks.NewStorage(t)
	store.On("BannerOff", mock.Anything, 1, 1).Return(nil)

	app := app.NewApp(store, nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	require.NoError(t, app.BannerOff(ctx, 1, 1))
}

func TestBanner(t *testing.T) {
	stats := storage.Stats{
		{ID: 1, Views: 9, Clicks: 7},
		{ID: 2, Views: 10, Clicks: 6},
		{ID: 3, Views: 12, Clicks: 6},
		{ID: 4, Views: 19, Clicks: 9},
		{ID: 5, Views: 24, Clicks: 11},
	}

	store := mocks.NewStorage(t)
	store.On("Stats", mock.Anything, 1, 1).Return(stats, nil)
	store.On("CountView", mock.Anything, 1, 1, 1).Return(nil)

	ntf := mocks.NewNotifier(t)
	ntf.On("Notify", mock.Anything, mock.Anything).Return(nil)

	app := app.NewApp(store, ntf)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	id, err := app.Banner(ctx, 1, 1)
	require.NoError(t, err)
	require.Equal(t, 1, id)
}
