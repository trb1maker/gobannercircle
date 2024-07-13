package cmd

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var migrationsFS embed.FS

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Выполнение миграций",
	Run: func(cmd *cobra.Command, _ []string) {
		os.Exit(migrate(cmd.Context()))
	},
}

func SendMigrations(fs embed.FS) {
	migrationsFS = fs
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func migrate(ctx context.Context) int {
	conn, err := sql.Open(
		"pgx",
		fmt.Sprintf(
			"host=%s port=%d dbname=%s user=%s password=%s",
			viper.GetString("storage.host"),
			viper.GetUint16("storage.port"),
			viper.GetString("storage.dbname"),
			viper.GetString("storage.user"),
			viper.GetString("storage.password"),
		),
	)
	if err != nil {
		slog.Error("can't init connection to postgres", "err", err)
		return failCode
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err = conn.PingContext(ctx); err != nil {
		slog.Error("can't connect to postgres", "err", err)
		return failCode
	}

	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("postgres"); err != nil {
		slog.Error("can't set migrations dialect", "err", err)
		return failCode
	}

	if err := goose.UpContext(ctx, conn, "migrations/postgres"); err != nil {
		slog.Error("can't execute migrations", "err", err)
		return failCode
	}

	return okCode
}
