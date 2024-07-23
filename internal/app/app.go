package app

import (
	"context"
	"math"
	"sort"
	"time"

	"github.com/trb1maker/gobannercircle/internal/notify"
	"github.com/trb1maker/gobannercircle/internal/storage"
)

func NewApp(storage Storage, notifier Notifier) *App {
	return &App{
		storage:  storage,
		notifier: notifier,
	}
}

type Storage interface {
	BannerOn(ctx context.Context, slotID, bannerID int) error
	BannerOff(ctx context.Context, slotID, bannerID int) error
	Stats(ctx context.Context, slotID, groupID int) (storage.Stats, error)
	IncViewCount(ctx context.Context, slotID, bannerID, groupID int) error
	IncClickCount(ctx context.Context, slotID, bannerID, groupID int) error
}

type Notifier interface {
	Notify(ctx context.Context, message notify.Message) error
}

type App struct {
	storage  Storage
	notifier Notifier
}

func (a *App) BannerOn(ctx context.Context, slotID, bannerID int) error {
	return a.storage.BannerOn(ctx, slotID, bannerID)
}

func (a *App) BannerOff(ctx context.Context, slotID, bannerID int) error {
	return a.storage.BannerOff(ctx, slotID, bannerID)
}

func (a *App) Banner(ctx context.Context, slotID, groupID int) (bannerID int, err error) {
	stats, err := a.storage.Stats(ctx, slotID, groupID)
	if err != nil {
		return 0, err
	}

	var pp, vv float64
	for i := 0; i < len(stats); i++ {
		pp += float64(stats[i].Clicks) / float64(stats[i].Views)
		vv += float64(stats[i].Views)
	}

	pp /= float64(len(stats))
	vv = 2 * math.Log(vv)

	for i := 0; i < len(stats); i++ {
		stats[i].P = pp + math.Sqrt(vv/float64(stats[i].Views))
	}

	sort.Stable(stats)
	bannerID = stats[len(stats)-1].ID

	if err := a.storage.IncViewCount(ctx, slotID, bannerID, groupID); err != nil {
		return 0, err
	}

	if err := a.notifier.Notify(ctx, notify.Message{
		Type:     "view",
		SlotID:   slotID,
		BannerID: bannerID,
		GroupID:  groupID,
		Time:     time.Now(),
	}); err != nil {
		return 0, err
	}

	return bannerID, nil
}

func (a *App) Click(ctx context.Context, slotID, bannerID, groupID int) error {
	if err := a.storage.IncClickCount(ctx, slotID, bannerID, groupID); err != nil {
		return err
	}

	if err := a.notifier.Notify(ctx, notify.Message{
		Type:     "view",
		SlotID:   slotID,
		BannerID: bannerID,
		GroupID:  groupID,
		Time:     time.Now(),
	}); err != nil {
		return err
	}

	return nil
}
