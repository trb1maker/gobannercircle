package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trb1maker/gobannercircle/internal/app"
	notify "github.com/trb1maker/gobannercircle/internal/notify/kafka"
	"github.com/trb1maker/gobannercircle/internal/service"
	"github.com/trb1maker/gobannercircle/internal/storage/postgres"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Запуск сервиса",
	Run: func(cmd *cobra.Command, _ []string) {
		os.Exit(runService(cmd.Context()))
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func runService(ctx context.Context) int {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	storage, err := postgres.NewPostgresStorage(
		viper.GetString("storage.host"),
		viper.GetUint16("storage.port"),
		viper.GetString("storage.dbname"),
		viper.GetString("storage.user"),
		viper.GetString("storage.password"),
	)
	if err != nil {
		slog.Error("can't init connection to postgres", "err", err)
		return failCode
	}

	if err = storage.Connect(ctx); err != nil {
		slog.Error("can't connect to postgres", "err", err)
		return failCode
	}
	defer storage.Close()

	notifier := notify.NewKafkaNotify(
		viper.GetString("notify.host"),
		viper.GetInt("notify.port"),
		viper.GetString("notify.topic"),
		viper.GetInt("notify.partition"),
	)
	if err = notifier.Connect(ctx); err != nil {
		slog.Error("can't connect to kafka", "err", err)
		return failCode
	}
	defer notifier.Close()

	logic := service.NewService(
		app.NewApp(storage, notifier),
		viper.GetString("service.host"),
		viper.GetUint16("service.port"),
	)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		logic.Stop(ctx)
	}()

	if err = logic.Start(); err != nil {
		slog.Error("service", "err", err)
		return failCode
	}

	return okCode
}
