package cmd

import (
	"context"
	"embed"
	"time"

	"github.com/jackc/pgx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var migrations embed.FS

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Выполнение миграций",
	Run: func(cmd *cobra.Command, _ []string) {
		conn, err := pgx.Connect(pgx.ConnConfig{
			Host:     viper.GetString("storage.host"),
			Port:     viper.GetUint16("storage.port"),
			Database: viper.GetString("storage.dbname"),
			User:     viper.GetString("storage.user"),
			Password: viper.GetString("storage.password"),
		})
		cobra.CheckErr(err)
		defer conn.Close()

		ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
		defer cancel()

		cobra.CheckErr(conn.Ping(ctx))

		migrationScripts, err := migrations.ReadDir("migrations/postgres")
		cobra.CheckErr(err)

		tx, err := conn.BeginEx(ctx, nil)
		cobra.CheckErr(err)
		defer tx.Rollback()

		for _, migrationScript := range migrationScripts {
			script, err := migrations.ReadFile(migrationScript.Name())
			cobra.CheckErr(err)

			_, err = tx.ExecEx(ctx, string(script), nil)
			cobra.CheckErr(err)
		}

		tx.CommitEx(ctx)
	},
}

func MigrationsPipe(m embed.FS) {
	migrations = m
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
