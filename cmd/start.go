package cmd

import (
	"context"
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
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
		defer cancel()

		storage, err := postgres.NewPostgresStorage(
			viper.GetString("storage.host"),
			viper.GetUint16("storage.port"),
			viper.GetString("storage.dbname"),
			viper.GetString("storage.user"),
			viper.GetString("storage.password"),
		)
		cobra.CheckErr(err)

		cobra.CheckErr(storage.Connect(ctx))
		defer storage.Close()

		notifier := notify.NewKafkaNotify(
			viper.GetString("logger.host"),
			viper.GetInt("logger.port"),
			viper.GetString("logger.topic"),
			viper.GetInt("logger.partition"),
		)
		cobra.CheckErr(notifier.Connect(ctx))
		defer notifier.Close()

		logic := service.NewService(
			app.NewApp(storage, notifier),
			viper.GetString("service.host"),
			viper.GetUint16("service.port"),
		)

		go func() {
			<-ctx.Done()

			ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
			defer cancel()
			logic.Stop(ctx)
		}()

		cobra.CheckErr(logic.Start())
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
